package main

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "time"
    "bufio"
    "os/exec"    

    "net/url" 
    "database/sql"
    _ "github.com/lib/pq"  
    "log" 
    "strings"
    "errors"
)

func logging_var(input_var interface{}) {

	//string msg
	//msg = fmt.Sprintf("%#v", input_var)

	var msg string

	msg = fmt.Sprintf("%#v", input_var)

	//this working good
	//fmt.Printf("%#v", input_var)

	//this not working
	//msg = "test:logging_var:inside"

	f2, err := os.OpenFile("/home/michael/Log/local_testing_golang_testnet.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		//log.Println(err)
	}
	defer f2.Close()

	if _, err := f2.WriteString(time.Now().String()); err != nil {
		//log.Println(err)
	}

	if _, err := f2.WriteString("|" + msg + "\n"); err != nil {
		//log.Println(err)
	}

}

func get_account_balance () {

	response, err := http.Get("http://127.0.0.1:9000/v2/accounts/ST1NBF7F98E6A42CYQKZA694ABW65T7JVX716WPP3")
    if err != nil {
        fmt.Printf("The HTTP request failed with error %s\n", err)
        logging_var (err)

    } else {
        data, _ := ioutil.ReadAll(response.Body)
        fmt.Println(string(data))
        logging_var (data)
    }

}



const (
  host     = "localhost"
  port     = 5432
  user     = "blockstack_testing_user_1"
  password = "password"
  dbname   = "blockstack_testing"  
)

func blockchain_run_create_node_testing () {

    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)

    db, err := sql.Open("postgres", psqlInfo)

    rows, err := db.Query("SELECT * FROM test_scenario_create_node")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    //println(rows)

    var id int 
    var root_dir string 
    var node_id string
    var scenario_description sql.NullString                
    var run_testing_flag sql.NullString    
    var run_testing_result_verbiage sql.NullString
    var comment sql.NullString 

    for rows.Next() {
        
        if err := rows.Scan(&id, &root_dir, &node_id, &scenario_description, &run_testing_flag, &run_testing_result_verbiage, &comment); err != nil {
            log.Fatal(err)
        }
        fmt.Printf("id = %n | root_dir=%s| node_id= %s |scenario_description=%s |run_testing_flag=%s|run_testing_result_verbiage=%s|comment=%s ", id,root_dir, node_id,scenario_description,run_testing_flag,run_testing_result_verbiage,comment)
    } 

    //blockchain_launch_node (root_dir, node_id)
    blockchain_launch_node (root_dir, node_id)
}

func search_on_log_file (logfilefulpath string, searchverbiagae string) {

    f, err := os.Open(logfilefulpath)
    if err != nil {
        //return 0, err
        fmt.Println("file operation error:")        
    }
    defer f.Close()

    // Splits on newlines by default.
    scanner := bufio.NewScanner(f)

    var line int64 
    line = 1
    // https://golang.org/pkg/bufio/#Scanner.Scan
    for scanner.Scan() {
        if strings.Contains(scanner.Text(), searchverbiagae) {
            //return line, nil
            fmt.Println("text:" + scanner.Text())        
            fmt.Println("line number:%d", line)
        }

        line++
    }

    if err := scanner.Err(); err != nil {
        // Handle the error
    }
}

func fetch_current_lineNum_log_file (logfilefulpath string) (int64) {

    f, err := os.Open(logfilefulpath)
    if err != nil {
        //return 0, err
        fmt.Println("file operation error:")        
    }
    defer f.Close()

    // Splits on newlines by default.
    scanner := bufio.NewScanner(f)

    var line int64 
    line = 1
    // https://golang.org/pkg/bufio/#Scanner.Scan
    for scanner.Scan() {        

        line++
    }

    if err := scanner.Err(); err != nil {
        // Handle the error        
    }

    return line
}

func fetch_current_section_log_file (logfilefulpath string, lastLineNumber int64) {

    f, err := os.Open(logfilefulpath)
    if err != nil {
        //return 0, err
        fmt.Println("file operation error:")        
    }
    defer f.Close()

    // Splits on newlines by default.
    scanner := bufio.NewScanner(f)

    var line int64 
    line = 1
    // https://golang.org/pkg/bufio/#Scanner.Scan
    for scanner.Scan() { 

         if (line >= lastLineNumber){

            fmt.Println("text:" + scanner.Text())        
            fmt.Println("line number:%d", line)       

         }

        line++
    }

    if err := scanner.Err(); err != nil {
        // Handle the error
    }
}

//cargo testnet start --config=./testnet/stacks-node/conf/neon-miner-conf.toml
func blockchain_launch_node (rootDir string, testingNodeID string) {

    var currentDir string

    path, err := os.Getwd()
    if err != nil {
        fmt.Println(err)
    }

    fmt.Println(path)     

    currentDir = rootDir + "/" + testingNodeID + "/stacks-blockchain"
    err = os.Chdir(currentDir)

    path, err = os.Getwd()
    fmt.Println("path1=" + path)    

    run_node_arg0 := "cargo" 
    run_node_arg1 := "testnet"
    run_node_arg2 := "start"
    run_node_arg3 := "--config=./testnet/stacks-node/conf/neon-miner-conf.toml"
    run_node_log_arg_1 :=">"
    run_node_log_arg_2 :="node_running_logging.txt"
    run_node_cmd := exec.Command(run_node_arg0, run_node_arg1, run_node_arg2, run_node_arg3,run_node_log_arg_1,run_node_log_arg_2)
    run_node_cmd.Run()
    run_node_stdout, run_node_err := run_node_cmd.Output()    
    if run_node_err != nil {
        fmt.Println(run_node_err.Error())        
    }
    fmt.Print("run_node_stdout =" + string(run_node_stdout))
}

func blockchain_run_testing () {

    fmt.Println("blockchain_run_testing:started")                     
 
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)

    db, err := sql.Open("postgres", psqlInfo)

    rows, err := db.Query("SELECT * FROM test_scenario_unit_testing")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    //println(rows)

    var id int 
    var test_scenario_id string 
    var test_scenario_input_parameters string
    var scenario_description sql.NullString                
    var run_testing_flag string
    var comment sql.NullString 

    for rows.Next() {
        
        if err := rows.Scan(&id, &test_scenario_id, &test_scenario_input_parameters, &scenario_description, &run_testing_flag, &comment); err != nil {
            log.Fatal(err)
        }
        //fmt.Println("fetching test_scenario_unit_testing table")
        //fmt.Printf("id = %n | test_scenario_id=%s| test_scenario_input_parameters= %s |scenario_description=%s |run_testing_flag=%s|comment=%s ", id,test_scenario_id, test_scenario_input_parameters,scenario_description,comment)

        if (run_testing_flag == "true") {

            db, err := sql.Open("postgres", psqlInfo)

            //fmt.Println("fetching testing_environment table")

            rows, err := db.Query("select testing_environment.test_working_directory as test_working_directory, testing_environment.rpc_port as rpc_port from testing_environment inner join running_scenario_on_testing_environment ON running_scenario_on_testing_environment.test_scenario_id = $1 " , test_scenario_id)
            if err != nil {
                log.Fatal(err)
            }
            defer rows.Close()
            //println(rows)

            var test_working_directory string 
            var rpc_port string
            
            for rows.Next() {
                
                if err := rows.Scan(&test_working_directory, &rpc_port); err != nil {
                    log.Fatal(err)
                }
                //fmt.Printf("test_working_directory = %n | rpc_port=%s ", test_working_directory,rpc_port)
            }    

            blockchain_command_call_post_API (test_working_directory, test_scenario_input_parameters, rpc_port)
        }
    } 
}

func blockchain_command_call_post_API (workingFolder string, cmdstr string, rpc_port string) {

    var currentLineNumber int64

    path, err := os.Getwd()
    if err != nil {
        fmt.Println(err)
    }

    err = os.Chdir(workingFolder)
    path, err = os.Getwd()    
    if err != nil {
        fmt.Println(err)
    }

    fmt.Println(path)

    cmdSlice := strings.Split(cmdstr, " ")
    fmt.Println("run_cmd=")        
    fmt.Println(cmdSlice)

    var arg [20]string
    var i int

    for i=0;i<20;i++ {
        arg[i] = ""
    }

    fmt.Println("cmdSlice")
    for i=0;i<len(cmdSlice);i++ {
        arg[i] = cmdSlice[i]            
    }  
    
    switch len(cmdSlice){
    case 9:
        run_cmd := exec.Command(arg[0],arg[1],arg[2],arg[3],arg[4],arg[5],arg[6],arg[7],arg[8])
        run_cmd_stdout, run_cmd_err := run_cmd.Output()    
        if run_cmd_err != nil {
        fmt.Println(run_cmd_err.Error())        
        }
        fmt.Print("run_cmd_stdout =" + string(run_cmd_stdout)) 
    case 11:
        run_cmd := exec.Command(arg[0],arg[1],arg[2],arg[3],arg[4],arg[5],arg[6],arg[7],arg[8],arg[9],arg[10])
        run_cmd_stdout, run_cmd_err := run_cmd.Output()    
        if run_cmd_err != nil {
        fmt.Println(run_cmd_err.Error())        
        }
        fmt.Print("run_cmd_stdout =" + string(run_cmd_stdout))  

        currentLineNumber = fetch_current_lineNum_log_file("/home/michael/blockstack_testNet/blockchain_launch_node/node1/stacks-blockchain/trace_logging_2020-05-25.txt")
        golang_api_post_deploy_tx(string(run_cmd_stdout), rpc_port)
        fetch_current_section_log_file("/home/michael/blockstack_testNet/blockchain_launch_node/node1/stacks-blockchain/trace_logging_2020-05-25.txt", currentLineNumber)
        
    case 14:
        run_cmd := exec.Command(arg[0],arg[1],arg[2],arg[3],arg[4],arg[5],arg[6],arg[7],arg[8],arg[9],arg[10],arg[11],arg[12],arg[13])
        run_cmd_stdout, run_cmd_err := run_cmd.Output()    
        if run_cmd_err != nil {
        fmt.Println(run_cmd_err.Error())        
        }
        fmt.Print("run_cmd_stdout =" + string(run_cmd_stdout))  

        currentLineNumber = fetch_current_lineNum_log_file("/home/michael/blockstack_testNet/blockchain_launch_node/node1/stacks-blockchain/trace_logging_2020-05-25.txt")
        golang_api_post_deploy_tx(string(run_cmd_stdout), rpc_port)
        fetch_current_section_log_file("/home/michael/blockstack_testNet/blockchain_launch_node/node1/stacks-blockchain/trace_logging_2020-05-25.txt", currentLineNumber)

        
    case 16:
        run_cmd := exec.Command(arg[0],arg[1],arg[2],arg[3],arg[4],arg[5],arg[6],arg[7],arg[8],arg[9],arg[10],arg[11],arg[12],arg[13],arg[14],arg[15])
        run_cmd_stdout, run_cmd_err := run_cmd.Output()    
        if run_cmd_err != nil {
        fmt.Println(run_cmd_err.Error())        
        }
        fmt.Print("run_cmd_stdout =" + string(run_cmd_stdout))  

        currentLineNumber = fetch_current_lineNum_log_file("/home/michael/blockstack_testNet/blockchain_launch_node/node1/stacks-blockchain/trace_logging_2020-05-25.txt")
        golang_api_post_deploy_tx(string(run_cmd_stdout), rpc_port)
        fetch_current_section_log_file("/home/michael/blockstack_testNet/blockchain_launch_node/node1/stacks-blockchain/trace_logging_2020-05-25.txt", currentLineNumber)

    case 17:
        run_cmd := exec.Command(arg[0],arg[1],arg[2],arg[3],arg[4],arg[5],arg[6],arg[7],arg[8],arg[9],arg[10],arg[11],arg[12],arg[13],arg[14],arg[15],arg[16])
        run_cmd_stdout, run_cmd_err := run_cmd.Output()    
        if run_cmd_err != nil {
        fmt.Println(run_cmd_err.Error())        
        }
        fmt.Print("run_cmd_stdout =" + string(run_cmd_stdout))    
    default:
        fmt.Println("not run, because parameter number not matched")                     
    }   
}

func golang_api_post_deploy_tx (binaryStr string, rpcPort string) error{

    print(len(binaryStr))

    var size int64 = int64( len(binaryStr) )

    var outputBytes = make([]byte, size)
    copy(outputBytes, binaryStr)

    rpc_url := "http://localhost:" + rpcPort + "/v2/transactions"

    resp, err := http.Post(rpc_url,
        "application/octet-stream", bytes.NewBuffer(outputBytes))

    if err != nil {
        print(err)
        return errors.New("custom error")
    }
    //defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        print(err)
    }

    fmt.Println("golang_api_post_deploy_tx:fun:respbody")

    fmt.Println(string(body))

    return nil
}

//1. gitclone, change oeml, start node, running by its own thread
func blockchain_setup_node (rootDir string, testingNodeID string) {

    var currentDir string

    path, err := os.Getwd()
    if err != nil {
        fmt.Println(err)
    }

    fmt.Println(path)     

    currentDir = rootDir
    err = os.Chdir(currentDir)

    path, err = os.Getwd()
    fmt.Println("path1=" + path)    

    arg0 := "mkdir" 
    arg1 := testingNodeID
    cmd := exec.Command(arg0, arg1)
    cmd.Run()
    stdout, err := cmd.Output()
    if err != nil {
        fmt.Println(err.Error())        
    }
    fmt.Print(string(stdout))

    currentDir = currentDir + "/" + testingNodeID
    err = os.Chdir(currentDir)
    path, err = os.Getwd()
    fmt.Println("path2=" + path)    
     
    //git clone https://github.com/blockstack/stacks-blockchain.git
    //git clone https://github.com/blockstack/stacks-blockchain.git >  logging-2020-05-12_2.txt
    gitclone_arg0 := "git" 
    gitclone_arg1 := "clone"
    gitclone_arg2 := "https://github.com/blockstack/stacks-blockchain.git"
    //gitclone_log_arg_1 :=">"
    //gitclone_log_arg_2 :="gitclone_logging.txt"
    
    gitclone_cmd := exec.Command(gitclone_arg0, gitclone_arg1, gitclone_arg2)
    gitclone_cmd.Run()
    gitclone_stdout, gitclone_err := gitclone_cmd.Output()
    if gitclone_err != nil {
        fmt.Println(gitclone_err.Error())        
    }
    fmt.Print("gitclone_stdout=" + string(gitclone_stdout))
    
    currentDir = currentDir + "/stacks-blockchain"
    err = os.Chdir(currentDir)

    //npx blockstack-cli@1.1.0-beta.1 make_keychain -t
    //npx blockstack-cli@1.1.0-beta.1 make_keychain -t |& tee -a  make_keychain_logging.txt
    npx_arg0 := "npx" 
    npx_arg1 := "blockstack-cli@1.1.0-beta.1"
    npx_arg2 := "make_keychain"
    npx_arg3 := "-t"
    npx_log_arg_1 :=">"
    npx_log_arg_2 :="make_keychain_logging.txt"
    npx_cmd := exec.Command(npx_arg0, npx_arg1, npx_arg2, npx_arg3,npx_log_arg_1,npx_log_arg_2)
    npx_cmd.Run()
    npx_stdout, npx_err := npx_cmd.Output()    
    if npx_err != nil {
        fmt.Println(npx_err.Error())        
    }
    fmt.Print("npx_stdout =" + string(npx_stdout))
    
}

func submit_BTC_request_faucet (){

    
    btcAddress := "n1HU9XDJynmvi8M3XsQQaYik7U4PpAEbMD"
    resp, err := http.PostForm("https://sidecar.staging.blockstack.xyz/sidecar/v1/faucets/btc",
    url.Values{"address": {btcAddress}, "id": {"submit"}})

    if err != nil {
        fmt.Printf("err:%v\n", err)
    }
    fmt.Printf("%+v", resp)
}

func main() {

    argsWithProg := os.Args    

    arg := os.Args

    fmt.Println("len of arg:%d" , len(arg) )    
        
    if (len(arg) != 2) {
        fmt.Println("help menu: ....")        
        return 
    }

    fmt.Println(argsWithProg)    
    fmt.Println(arg[1])
    
    switch arg[1] {
    case "1":
        search_on_log_file("/home/michael/blockstack_testNet/blockchain_launch_node/node1/stacks-blockchain/trace_logging-2020-05-23_2-bk.txt", "Failed to parse HTTP request")    
        return
    
    case "2":
        blockchain_run_testing()
        return
    default:
        fmt.Println("wrong input flag")
        return 
        }  
}
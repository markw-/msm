// Multi-layer Synthetic Monitor (MSM) - Main package
//  
//
// Copyright (c) 2016 Mark Broughton Wild. All Rights Reserved.

package main

import (
    "fmt"
    "os"
    "io/ioutil"
    "strings"
    "regexp"
    "net/http"
    "net"
    "time"
    "bufio"
    "strconv"
    "bytes"
)

type check struct {
    check_type  string
    mode   		string
    target 		string
    text   		string
    uid    		string
    pw     		string
    result		string
}

func main(){
	cwd, _ := os.Getwd()
    fmt.Println("Starting Multi-layer Synthetic Monitor\nWorking directory -", cwd ,"\n")
    var check_period = 300
    var check_timeout = 1000
    var check_list []check
    var do_check_count int
    check_response := make(chan check)
    if len(os.Args) == 2 {
        check_list = parse_checks(os.Args[1])
    } else{
        check_list = parse_checks("msm_checks.txt")
    }
    for _, c := range check_list{
    	do_check_count++
    	// fmt.Printf("Check counter is incremented to %d\n", do_check_count)
    	go do_check(c, check_response)
    }
    fmt.Printf("Check counter value is %d\n", do_check_count)
    for {
    	if do_check_count == 0{
    		close(check_response)
    		fmt.Println("All checks finished. Exiting.")
    		os.Exit(0)
    	} else{
    		response := <- check_response
    		fmt.Printf("Check type - %s, mode - %s, target - %s, text - %s, uid - %s.\nResult: %s\n" , response.check_type, response.mode, response.target, response.text, response.uid, response.result)
    		do_check_count--
    		fmt.Printf("Check counter is decremented to %d\n", do_check_count)
    	}
    }
}

// This function parses the configuration file, reading each line to find the 
// program variables. The program variables are set directly by this parsing
// function (or suitable defaults are used) as there doesn't seem to be any
// good reason to pass them back to the main goroutine.
func parse_conf (filename string)

    re := regexp.MustCompile("(check_period|check_timeout)\\s*=\\s*\\d{2,}")
    
    // Open the config file
    conf_file, err := ioutil.ReadFile(filename)
    if err != nil{
        fmt.Printf("Can't read the config file %s\n%s\nUsing default settings.", filename, err)
        return
    }
    // First split out the lines and, for each line, check if it is a comment line (in which case ignore it)
    for _, line := range strings.Split(string(conf_file), "\n") {
    	if strings.HasPrefix(line, "//"){
    		//fmt.Printf("Comment line, skipping - %s\n", line)
    		continue
    	}
    	// The line is not a comment line so look for KV pairs. If no KV pairs are found, say so.
    	//fmt.Println(reflect.TypeOf(line),"\t-",line)
    	keyvaluepairs := re.FindAllString(line, -1)
    	if keyvaluepairs == nil{
    		//fmt.Printf("No key/value pairs in the line %q\n", line)
    		continue
    	} else if {
    		
        

// This function parses the checks file, reading each line and checking
// the values to ensure that they are valid for the check mode. Then writing
// the check values into a struct of type "check" and finally returning a slice
// of all the checks.

func parse_checks(filename string) []check {
	//{{{
    // Declare a slice of check structs and a regex pattern type now as we will possibly use them throughout the rest
    // of the function.
	var checks []check
    //re := regexp.MustCompile("(type|result|target|text|UID|PW)=((\\d|\\.|:)+|\\w+|\".+?\")")
    re := regexp.MustCompile("(type|result|target|text|UID|PW)=((\\w|:|/|\\.)+|((\\d{1,3}\\.){3}\\d{1,3}:\\d{5})|\\w+|\".+?\")")
    
    //re := regexp.MustCompile("check=\\w*")
    // Open the checks file (useful!)
    checks_file, err := ioutil.ReadFile(filename)
    if err != nil{
        fmt.Printf("Can't read the checks file %s\n%s\nCannot continue, exiting.", filename, err)
        os.Exit(1)
    }
    // First split out the lines and, for each line, check if it is a comment line (in which case ignore it)
    for _, line := range strings.Split(string(checks_file), "\n") {
    	if strings.HasPrefix(line, "//"){
    		//fmt.Printf("Comment line, skipping - %s\n", line)
    		continue
    	}
    	// If the line is not a comment line, look for KV pairs.
    	//fmt.Println(reflect.TypeOf(line),"\t-",line)
    	keyvaluepairs := re.FindAllString(line, -1)
    	if keyvaluepairs == nil{
    		//fmt.Printf("No key/value pairs in the line %q\n", line)
    		continue
    	}
        // Now find the keys and values for the checks from the line and populate a check struct with them
        var new_check check
        for _ , keyvaluepair := range keyvaluepairs{
        	//fmt.Println("Checking key\\value pair to populate the check type.") 
            key_value:= strings.Split(keyvaluepair, "=")
            key, value := key_value[0], key_value[1]
            //fmt.Printf("Key=%s, Value=%s\n", key, value)
            key=strings.ToLower(key)
            switch key{
            case "type":
            	value = strings.ToLower(value)
                switch value{
                case "application":
                    new_check.check_type = "app"
                case "url":
                    new_check.check_type = "url"
                case "db":
                    new_check.check_type = "db"
                }
            
            case "result":
            	value = strings.ToLower(value)
                switch value{
                case "alive":
                    new_check.mode = "alive"
                case "timed":
                    new_check.mode = "timed"
                }
                
            case "target":
                new_check.target = strings.Trim(value, "\"")
              
            case "text":
                new_check.text = strings.Trim(value, "\"")

            case "uid":
                new_check.uid = strings.Trim(value, "\"")
                
            case "pw":
                new_check.pw = strings.Trim(value, "\"")
                
            default:
                fmt.Printf("Unknown key %s", key)
            }           
        }
        // Add the new check struct to the slice of checks so it is in the type returned by the function
        checks = append(checks, new_check)
    }
    return checks
}
//}}}

// Function to be called in a new goroutine, takes a check struct, carries out 
// the check specified in that struct and writes the result to stdout/syslog/
// the results file. It calls the function appropriate to the check being made
// (application, URL or DB).

func do_check (c check, response_chan chan <- check){
	//{{{
	if c.check_type == "app"{
		result, _, _ := check_app(c.mode, c.target)
		if result == true {
			c.result ="Target is ALIVE."
		} else {
			c.result ="Target is DEAD."
		}
	} else if c.check_type == "url"{
		result, response_time := check_url(c.mode, c.target)
		if result == true{
			if c.mode == "alive"{
				c.result ="Target is ALIVE."
			}
			if c.mode == "timed"{
				var buffer bytes.Buffer
				buffer.WriteString("Target responded in ")
				buffer.WriteString(strconv.Itoa(response_time))
				buffer.WriteString(" miliseconds.")
				c.result = buffer.String()
			}
		} else {
			if c.mode == "alive"{
				c.result ="Target is DEAD."
			}
			if c.mode == "timed"{
				c.result ="Target did not respond."
			}
			
		}	
	} else if c.check_type == "db"{
		result, db_id, _ := check_db(c.mode, c.target)
		if result == true {
			c.result ="Target is ALIVE. DB identified as " +db_id+"."
		} else {
			c.result ="Target is DEAD."
		}
	} else{
		c.result =" ERROR - Unknown check type."
	}
	
	response_chan <- c
}
//}}}

// This function carries out an application check

func check_app(mode string, address string, timeout time.Duration) (response bool, application_id string, response_time int){
	//{{{
	//fmt.Printf("Checking Application - %s | type of check - %s", adress, mode)
	response = false
	response_time = -1
	application_id = "Unknown application"
	const id_len = 40
	var id [id_len]byte
	
	// Connect to the application target (only IP and port at the moment, need to add options for names/ODBC/etc.)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return
	}
	defer conn.Close()
	id_reader := bufio.NewReader(conn)
	if err != nil {
		return
	}
	i := id_len
	for i < id_len {
		id[i], err = id_reader.ReadByte()
		if err != nil{
			break
		}
		i++
	}
	
	// ID the application if possible or just note that an application respsonded
    if strings.Contains(string(id[:id_len]), "RFB"){
		application_id = "VNC (Remote Frame Buffer)"
		response = true
	} else {
		response = true
	}
	
	return
}
//}}}

// This function carries out a URL check

func check_url(mode string, url string) (response bool, response_time int) {
	//{{{ 
	//fmt.Printf("Checking URL - %s | type of check - %s", url, mode)
	response = false
	response_time = -1
	if mode == "alive"{
		_ , err := http.Get(url)
		if err != nil{
			return
			}
		response = true
		}
		
	if mode == "timed"{
		t0 := time.Now()
		_ , err := http.Get(url)
		t1:= time.Now()
		if err != nil{
			return
		}
		precise := t1.Sub(t0)
		duration := int(precise/time.Millisecond)
		response = true
		response_time = duration
	}
	return
}
//}}} 

// This function carries out a DB check

func check_db(mode string, db_address string, timeout time.Duration) (response bool, db_id string, response_time int){
	//{{{ 
	//fmt.Printf("Checking DB - %s | mode - %s", db_address, mode)
	response = false
	response_time = -1
	db_id = "Unknown DB"
	// Connect to the DB target (only IP and port at the moment, need to add options for names/ODBC/etc.)
	conn, err := net.DialTimeout("tcp", db_address, timeout)
	if err != nil {
		return
	}
	defer conn.Close()
	
	// Now try to ID the database or just note that there was a DB there
	id, err := bufio.NewReader(conn).Peek(40)
	if strings.Contains(string(id), "MariaDB"){
		db_id = "MariaDB"
		response = true
	} else {
		response = true
	}

	return
}
//}}}
// Package main implements a proof-of-concept for a test automation
// server (TAS).
//
// The idea is simple: a central server maintains a to-do list for test
// executors (TXs).  TXs and user-facing clients communicate with the TAS
// via a REST interface.  This TAS provides no authentication; a TX or
// user-client may invoke any REST end-point it chooses.
//
// This TAS is a spike, both to learn Go and to learn
// more about the relationship between TAS, TX, and user client.
//
// Eventually, this code will be thrown out and a real TAS/TX/client
// architecture documented, constructed, and tested.
package main

import (
  "bytes"
  "fmt"
  "log"
  "strings"
  "strconv"
  "net/http"
)

const (
  ok = "OK"
  bad = "BAD"
)

// jobQueue provides a mapping from job ID to some arbitrary blob of data,
// expressed here as a string.  At this point, the TAS doesn't care what
// is inside the blob.  It only cares about the job ID.
var jobQueue map[int64] string;

// nextId provides a unique ID for jobs.  Being 63 bits in size, and
// assuming one job creation per microsecond consistently, we expect
// roll-over problems to occur roughly 293,000 years in the future
// from the point of TAS start-up.  If you allow job IDs to go negative,
// we can extend this range to over 580,000 years.
var nextId int64 = 0;


func init() {
  jobQueue = make(map[int64] string);
}


// request logs that a request is being handled.  Calling this procedure
// is a considerate thing to do, but is not required by any handler.  Be
// considerate; when maintaining this server to add new end-points, please
// call this procedure so that logging remains consistent.
func request(method, path, status string) {
  log.Printf("Received %s request on %s (%s)", method, path, status)
}


// todoGet handles GET requests on the jobs/todo endpoint.  It produces a
// JSON array of jobs pending, identified by ID.
func todoGet(w http.ResponseWriter, r *http.Request) {
  request(r.Method, r.URL.Path, ok);
  fmt.Fprintf(w, "{\"pendingJobs\":[");
  for id := range jobQueue {
    fmt.Fprintf(w, "\"%d\"", id);
    if id != nextId-1 {
      fmt.Fprintf(w, ",");
    }
  }
  fmt.Fprintf(w, "]}");
}

// todoPost handles POST requests on the jobs/todo endpoint.  It takes an
// arbitrary body and associates it with a brand new job ID.  The result
// is the job ID passed back as a JSON construct if successful, or an
// error construct otherwise.
func todoPost(w http.ResponseWriter, r *http.Request) {
  var body bytes.Buffer
  var segment []byte

  request(r.Method, r.URL.Path, ok);

  segment = make([]byte, 4096);
  for n, err := r.Body.Read(segment); (err == nil) && (n >= 0); {
    body.Write(segment[0:n])
    n, err = r.Body.Read(segment);
  }
  jobQueue[nextId] = body.String()
  fmt.Fprintf(w, "{\"jobId\":\"%d\"}", nextId);
  nextId++;
}

// badRequest signals a 400 Bad Request response to the client.  The body
// is a descriptor of URL and method used.
func badRequest(w http.ResponseWriter, r *http.Request) {
  request(r.Method, r.URL.Path, bad);
  w.WriteHeader(400);
  fmt.Fprintf(w, "{\"error\": \"Illegal request\", \"method\": %q, \"url\": %q}", r.Method, r.URL.String());
}

// todoHandler dispatches requests based on GET vs POST methods to the 
// jobs/todo endpoint.
func todoHandler(w http.ResponseWriter, r *http.Request) {
  switch {
    case r.Method == "GET":  todoGet(w, r)
    case r.Method == "POST": todoPost(w, r)
    default: badRequest(w, r)
  }
}

// jobGet retrieves the arbitrary blob of data associated with a job.
// Typically used by a TX to figure out what test to run, or a
// client application to provide the user with a human-readable
// representation of the job.
//
// The job ID appears as a number on the end of the URI path; e.g.:
//
//    http://localhost:8080/job/14
//
// If the ID provided maps to a non-existing job, a successful
// result ensues, albeit with an empty response.
//
// If the ID provided cannot be used for some reason, a bad request
// error ensues.
func jobGet(w http.ResponseWriter, r *http.Request) {
  components := strings.Split(r.URL.Path, "/")
  id := components[len(components)-1];
  jobId, err := strconv.ParseInt(id, 10, 64)
  if err == nil {
    fmt.Fprintf(w, jobQueue[jobId]);
  } else {
    badRequest(w, r);
  }
}

// jobHandler dispatches to the GET, ... handlers for the
// job/{{id}} endpoint.
func jobHandler(w http.ResponseWriter, r *http.Request) {
  switch {
    case r.Method == "GET": jobGet(w, r)
    default: badRequest(w, r)
  }
}

// main configures and runs the web server responsible for handling
// TAS requests.
func main() {
  http.HandleFunc("/jobs/todo/", todoHandler);
  http.HandleFunc("/job/", jobHandler);
  log.Fatal(http.ListenAndServe(":8080", nil));
}


package main

import (
    "database/sql"
    "flag"
    "html/template"
	"log"
	"net/http"
    "os"

    "snippetbox.tegardp.com/internal/models"
    _ "github.com/go-sql-driver/mysql"
)

type application struct {
    errorLog *log.Logger
    infoLog *log.Logger
    snippets *models.SnippetModel
    templateCache map[string]*template.Template
}

func main() {
    // define command line parameters
	addr := flag.String("addr", ":4000", "HTTP network address")
    dsn := flag.String("dsn", "root:nirvana456@tcp(localhost:3306)/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()
    
    // create log format
    infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
    errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

    // connect db
    db, err := openDB(*dsn)
    if err != nil {
        errorLog.Fatal(err)
    }
   
    defer db.Close()

    templateCache, err := newTemplateCache()
    if err != nil {
        errorLog.Fatal(err)
    }

    app := &application{
        errorLog: errorLog,
        infoLog: infoLog,
        snippets: &models.SnippetModel{DB: db},
        templateCache: templateCache,
    }

    srv := &http.Server{
        Addr:       *addr,
        ErrorLog:   errorLog,
        Handler:    app.routes(),
    }

	infoLog.Println("Starting server on", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }
    if err = db.Ping(); err != nil {
        return nil, err
    }

    return db, nil
}

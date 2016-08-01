// Copyright (c) 2016 SeongJae Park.
//
// This program is free software; you can redistribute it and/or modify it
// under the terms of the GNU General Public License version 3 as published by
// the Free Software Foundation.

package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// document contains informations for a document.
type document struct {
	Name string
	Id   int // id of the document.
}

type documents_info struct {
	Docs    []document
	Next_id int
}

type wordinfo struct {
	Word         string
	Totalfreq    int
	Freq         map[string]int
	Succ_history []time.Time
	Fail_history []time.Time
}

type wordinfos struct {
	Wordinfos map[string]wordinfo
}

const (
	CMD          = "ago"
	USAGE        = "USAGE: " + CMD + " <commands> [argument ...]\n"
	NOARG_ERRMSG = USAGE + "\nFor detail, try " + CMD + " help\n"
	COMMANDS_MSG = `
commands are:
 - ls-docs, add-docs, rm-docs: list, add, remove documentation[s].
 - words: list words information
 - test: start a test. Number of questions can be specified as option.
 - doc, mod-doc: Commands for future. Not be implemented yet. Display and
                 modify information of the doc.

`
	HELP_MSG     = USAGE + COMMANDS_MSG
	ANDRD        = "android"
	ANDRD_TMPDIR = "/data/local/tmp"
	DOCSDIR      = "docs"
	DOCINFO      = "info"
	WORDINFO     = "words"
	DOCDIR_PREF  = "doc" // prefix of document directory
	DOCID_DICT   = 0     // document id for dictionary
	NR_RSVD_DOCS = 1
)

var (
	errl        = log.New(os.Stderr, "[err] ", 0)
	dbgl        = log.New(os.Stderr, "[dbg] ", 0)
	metadat_dir = "/tmp/.ago"
	docs_dir    string
	doci_path   string // documents information file path
	wordi_path  string // words information file path
	docs_info   documents_info
	winfos      wordinfos
)

func read_docs_info() {
	c, err := ioutil.ReadFile(doci_path)
	if err != nil {
		errl.Printf("failed to read doc info file: %s\n", err)
		os.Exit(1)
	}

	if err := json.Unmarshal(c, &docs_info); err != nil {
		errl.Printf("error while unmarshal doc info: %s\n", err)
		dbgl.Printf("the json: %s\n", c)
		os.Exit(1)
	}
}

func write_docs_info() {
	bytes, err := json.Marshal(docs_info)
	if err != nil {
		errl.Printf("failed to marshal docs_info: %s\n", err)
		os.Exit(1)
	}

	if err := ioutil.WriteFile(doci_path, bytes, 0600); err != nil {
		errl.Printf("failed to write marshaled docs_info: %s\n", err)
	}
}

func read_words_info() {
	c, err := ioutil.ReadFile(wordi_path)
	if err != nil {
		errl.Printf("failed to read words info file: %s\n", err)
		os.Exit(1)
	}

	if err := json.Unmarshal(c, &winfos); err != nil {
		errl.Printf("error while unmarshal words info: %s\n", err)
		dbgl.Printf("the json: %s\n", c)
		os.Exit(1)
	}
}

func write_words_info() {
	bytes, err := json.Marshal(winfos)
	if err != nil {
		errl.Printf("failed to marshal words_info: %s\n", err)
		os.Exit(1)
	}

	if err := ioutil.WriteFile(wordi_path, bytes, 0600); err != nil {
		errl.Printf("failed to write marshaled words_info: %s\n", err)
	}
}

func lsdocs(args []string) {
	for _, doc := range docs_info.Docs {
		fmt.Printf("%d: %s\n", doc.Id, doc.Name)
	}
}

// file_exists checks whether a file of specific path exists.
//
// Returns true if exists, false if not.
func file_exists(path string) bool {
	if _, err := os.Stat(docs_dir); err == nil {
		return true
	}
	return false
}

func add_word(word string, freq, docid int) {
	winfo, exists := winfos.Wordinfos[word]
	if !exists {
		winfo = wordinfo{Word: word}
		winfo.Freq = make(map[string]int)
		winfo.Totalfreq = 0
	}
	winfo.Totalfreq += freq
	winfo.Freq[strconv.Itoa(docid)] += freq
	winfos.Wordinfos[word] = winfo
}

func analyze_words(bytes []byte, docid int) {
	freq_map := make(map[string]int)
	s := string(bytes)
	words := strings.Fields(s)
	for _, word := range words {
		word = strings.ToLower(word)
		for _, special := range []string{",", ".", "!", "?",
			"`", "'", "\"", ":", "(", ")"} {
			word = strings.Join(strings.Split(word, special), "")
		}
		if word == "" {
			continue
		}
		freq, _ := freq_map[word]
		freq_map[word] = freq + 1
	}
	for word, freq := range freq_map {
		add_word(word, freq, docid)
	}
}

func adddoc(file_path string) error {
	if !file_exists(file_path) {
		err := errors.New("file not exists")
		return err
	}

	// read the file
	bytes, err := ioutil.ReadFile(file_path)
	if err != nil {
		msg := fmt.Sprintf("failed to read file: %s", err)
		err := errors.New(msg)
		return err
	}

	docid := docs_info.Next_id

	// analyze words in the file content
	analyze_words(bytes, docid)

	// create dir under docs/
	docdir := fmt.Sprintf("%s%d", DOCDIR_PREF, docid)
	docdirpath := path.Join(docs_dir, docdir)
	if err = os.MkdirAll(docdirpath, 0700); err != nil {
		msg := fmt.Sprintf("failed to create dir: %s", err)
		err = errors.New(msg)
		return err
	}

	// write copy in the doc<doc id>/
	_, docname := filepath.Split(file_path)
	in_file_path := path.Join(docdirpath, docname)
	if err = ioutil.WriteFile(in_file_path, bytes, 0600); err != nil {
		msg := fmt.Sprintf("failed to write file: %s", err)
		err = errors.New(msg)
		return err
	}

	// add to docs_info global object
	doc := document{Name: docname, Id: docid}
	docs_info.Next_id += 1
	docs_info.Docs = append(docs_info.Docs, doc)
	return nil
}

func adddocs(args []string) {
	for _, path := range args {
		if err := adddoc(path); err != nil {
			errl.Printf("failed to add doc %s: %s\n",
				path, err)
			os.Exit(1)
		}
	}
	write_docs_info()
	write_words_info()
}

// rmdoc remove a document with specific id.
//
// Returns true if success to remove the document, false if not
func rmdoc(target int) error {
	for idx, doc := range docs_info.Docs {
		if doc.Id != target {
			continue
		}
		docdir := fmt.Sprintf("%s%d", DOCDIR_PREF, doc.Id)
		docdirpath := path.Join(docs_dir, docdir)
		if err := os.RemoveAll(docdirpath); err != nil {
			msg := fmt.Sprintf("failed to remove dir: %s", err)
			err = errors.New(msg)
			return err
		}

		docs_info.Docs = append(docs_info.Docs[:idx],
			docs_info.Docs[idx+1:]...)
		return nil
	}
	return errors.New("no such doc")
}

func rmdocs(args []string) {
	for _, arg := range args {
		target, err := strconv.Atoi(arg)
		if err != nil {
			errl.Printf("argument must be doc id. err: \n", err)
		}
		if err = rmdoc(target); err != nil {
			fmt.Printf("failed to remove doc id %d: %s\n",
				target, err)
		}
	}

	write_docs_info()
}

type ordered_wis []wordinfo

func (a ordered_wis) Len() int {
	return len(a)
}

func (a ordered_wis) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ordered_wis) Less(i, j int) bool {
	ltl := a[i]
	big := a[j]

	lf := len(ltl.Fail_history)
	ls := len(ltl.Succ_history)
	bf := len(big.Fail_history)
	bs := len(big.Succ_history)

	l_imp := ltl.Totalfreq + 100*(lf-ls)
	b_imp := big.Totalfreq + 100*(bf-bs)

	// We need descendent order
	return l_imp > b_imp
}

func lswords(args []string) {
	var slices []wordinfo
	for _, info := range winfos.Wordinfos {
		slices = append(slices, info)
	}

	sort.Sort(ordered_wis(slices))
	for _, info := range slices {
		fmt.Printf("word %s:\t total freq %d, succ %d, fail %d\n",
			info.Word, info.Totalfreq,
			len(info.Succ_history), len(info.Fail_history))
	}
}

// html_to_txt extracts text that browser displays from html source.
func html_to_txt(s string) string {
	for {
		scriptopen := strings.Index(s, "<script")
		scriptclose := strings.Index(s, "</script>")
		if scriptopen == -1 || scriptclose == -1 {
			break
		}
		s = s[:scriptopen] + s[scriptclose+len("</script>"):]
	}
	for {
		tagopen := strings.IndexByte(s, '<')
		tagclose := strings.IndexByte(s, '>')
		if tagopen == -1 || tagclose == -1 || tagopen > len(s) ||
			tagclose > len(s) {
			break
		}
		if tagopen > tagclose {
			fmt.Printf("something wrong: %d %d %s\n",
				tagopen, tagclose, s[tagclose:tagclose+200])
			os.Exit(1)
			break
		}

		s = s[:tagopen] + s[tagclose+1:]
	}
	return s
}

// pick_txt picks up a section inside a text.
// the section should starts / ends with specific text, start and end.
// If the section not found, it just returns blank string.
func pick_section(s string, start string, end string) string {
	sidx := strings.Index(s, start)
	if sidx == -1 {
		return ""
	}
	eidx := sidx + strings.Index(s[sidx+1:], end)
	if eidx == -1 {
		return ""
	}

	return s[sidx:eidx]
}

func nth_index(s, sep string, n int) int {
	found := 0
	ret := 0
	for found < n {
		if ret != 0 {
			ret += 1
		}
		ret += strings.Index(s[ret:], sep)
		if ret == -1 {
			return -1
		}
		found++
	}
	return ret
}

const DAUM_WORD = "<div class=\"card_word\""

func mean_section(s string) string {
	return pick_section(s, DAUM_WORD, DAUM_WORD)
}

func ex_section(s string) string {
	the_index := nth_index(s, DAUM_WORD, 3)
	return pick_section(s[the_index:], DAUM_WORD, DAUM_WORD)
}

func excludeStrRegion(s, start, end string) string {
	sidx := strings.Index(s, start)
	if sidx == -1 {
		return ""
	}
	eidx := strings.Index(s, end)
	return s[:sidx] + s[eidx+utf8.RuneCountInString(end):]
}

type Node struct {
	XMLName xml.Name
	Content []byte `xml:",innerxml"`
	Class   []byte `xml:"class,attr"`
	Nodes   []Node `xml:",any"`
}

func get_node_with(nodes []Node, name, class string) *Node {
	for _, n := range nodes {
		if n.XMLName.Local == name && string(n.Class) == class {
			return &n
		}
		ret := get_node_with(n.Nodes, name, class)
		if ret != nil {
			return ret
		}
	}

	return nil
}

func spr_leaf_string(nodes []Node, buf *string) {
	for _, n := range nodes {
		if len(n.Nodes) == 0 {
			*buf += strings.TrimSpace(string(n.Content)) + " "
		}
		if n.XMLName.Local == "span" && string(n.Class) == "txt_ex" {
			*buf += "\n"
		}
		spr_leaf_string(n.Nodes, buf)
	}
}

func leaf_string(node Node) string {
	buf := ""
	spr_leaf_string([]Node{node}, &buf)
	return buf
}

func daum_dict(q string) string {
	daumdic_url := "http://dic.daum.net/search.do?q="
	suffix := "&dic=eng&search_first=Y"
	url := fmt.Sprintf("%s%s%s", daumdic_url, q, suffix)
	resp, err := http.Get(url)
	if err != nil {
		panic(fmt.Sprintf("error while get %s: %s", url, err))
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		panic(fmt.Sprintf("failed to read body: %s", err))
	}

	html_src := string(body)
	for {
		e := excludeStrRegion(html_src, "<script", "</script>")
		if e == "" {
			break
		}
		html_src = e
	}

	dec := xml.NewDecoder(bytes.NewBuffer([]byte(html_src)))
	dec.Strict = false
	dec.AutoClose = xml.HTMLAutoClose
	dec.Entity = xml.HTMLEntity

	var n Node
	err = dec.Decode(&n)
	if err != nil {
		panic("decode fail")
	}

	ret := ""

	relate_node := get_node_with([]Node{n}, "div", "card_relate")
	if relate_node != nil {
		ret += "# Related Words\n" + leaf_string(*relate_node) + "\n\n"
	}

	word_node := get_node_with([]Node{n}, "div", "card_word #word #word")
	if word_node != nil {
		ret += "# Word\n" + leaf_string(*word_node) + "\n\n"
	}

	mean_node := get_node_with([]Node{n}, "div", "card_word #word #mean")
	if mean_node != nil {
		ret += "# Meaning\n" + leaf_string(*mean_node) + "\n\n"
	}

	ex_node := get_node_with([]Node{n}, "div", "card_word")
	if ex_node != nil {
		ret += "# Examples\n" + leaf_string(*ex_node) + "\n\n"
	}

	return ret

	/*
		 XML parsing based current code is not complete yet.  It
		 ignores nodes like `<a src="abc">abc`.  For possible later
		 use, leave olde code here as comment for now.

		mean_sect := html_to_txt(mean_section(html_src))
		ex_sect := html_to_txt(ex_section(html_src))
		mean_sect = strings.Join(strings.Fields(mean_sect), " ")
		ex_sect_fields := strings.Fields(ex_sect)
		for idx, f := range ex_sect_fields {
			if strings.HasSuffix(f, ".") && len(f) > 2 {
				ex_sect_fields[idx] = f + "\n"
			}
		}
		ex_sect = strings.Join(ex_sect_fields, " ")
		return fmt.Sprintf("Meaning\n%s\n\nExamples\n%s\n",
			html_to_txt(mean_sect), html_to_txt(ex_sect))
	*/
}

func dic(args []string) {
	word := args[0]

	fmt.Printf("%s\n", daum_dict(args[0]))

	winfo, exists := winfos.Wordinfos[word]
	if !exists {
		add_word(word, 500, DOCID_DICT)
		winfo, _ = winfos.Wordinfos[word]
	}
	winfo.Fail_history = append(winfo.Fail_history, time.Now())
	winfos.Wordinfos[winfo.Word] = winfo
	write_words_info()
}

func get_questions(args []string) []wordinfo {
	pool := []wordinfo{}
	ret := []wordinfo{}
	wis := winfos.Wordinfos
	nr_ques := 5

	if len(args) > 0 {
		nr, err := strconv.Atoi(args[0])
		if err != nil {
			errl.Printf("wrong argument %s. it must be integer\n",
				args[0])
		} else {
			nr_ques = nr
		}
	}

	// Get random nr_ques * 10 words
	for _, wi := range wis {
		pool = append(pool, wi)
		if len(ret) == nr_ques*10 {
			break
		}
	}

	// Sort and use high score ones
	sort.Sort(ordered_wis(pool))
	if len(pool) < nr_ques {
		return pool
	}
	return pool[:nr_ques]
}

func do_singletest(wi wordinfo, date time.Time, nr, total int) {
	input := ""
	question := wi.Word

	fmt.Printf("==============================================\n")
	fmt.Printf("Question %d/%d:\n", nr, total)
	fmt.Printf("\n [[ %s ]]\n", wi.Word)
	fmt.Printf("\n\nPress Enter after you remember the meaning of it:\n")
	fmt.Printf("\n Reference: %d Test success/fail: %d/%d\n",
		wi.Totalfreq, len(wi.Succ_history), len(wi.Fail_history))
	fmt.Scanln(&input)
	fmt.Printf("----------------------------------------------\n")
	fmt.Printf("The maning of %s was:\n%s\n\n",
		question, daum_dict(question))
	fmt.Printf("----------------------------------------------\n")
	fmt.Printf("Were you understanding it well? (Yes/No)\n")
	fmt.Scanln(&input)
	if strings.HasPrefix(strings.ToUpper(input), "Y") {
		wi.Succ_history = append(wi.Succ_history, date)
		fmt.Printf("your feedback, Yes applied\n")
	} else {
		wi.Fail_history = append(wi.Fail_history, date)
		fmt.Printf("your feedback, No applied\n")
	}
	winfos.Wordinfos[wi.Word] = wi
	fmt.Printf("----------------------------------------------\n")
	fmt.Printf("\n\n\n")
}

func do_test(args []string) {
	fmt.Printf("Let the game begin\n\n")
	fmt.Printf("Ready? (Yes/[No])\n")
	input := ""
	fmt.Scanln(&input)
	if strings.HasPrefix(strings.ToUpper(input), "N") {
		fmt.Printf("OK, see you later ;)\n")
		return
	}

	questions := get_questions(args)
	date := time.Now()
	for idx, wi := range questions {
		do_singletest(wi, date, idx+1, len(questions))
		write_words_info()
	}
}

// main is the entry point of `ago`.
// ago usage is similar to familiar tools:
// 	ago <command> [argument ...]
//
// commands are:
// - ls-docs, add-docs, rm-docs: list, add, remove documentation[s].
// - words: list words information
// - test: start a test. Number of questions can be specified as option.
// - doc, mod-doc: Commands for future. Not be implemented yet. Display and
// 	modify information of the doc.
func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("No argument.\n")
		fmt.Printf(NOARG_ERRMSG)
		os.Exit(1)
	}
	cmd := args[1]
	args = args[2:]
	switch cmd {
	case "ls-docs":
		lsdocs(args)
	case "add-docs":
		adddocs(args)
	case "rm-docs":
		rmdocs(args)
	case "words":
		lswords(args)
	case "dic":
		dic(args)
	case "test":
		do_test(args)
	case "help":
		fmt.Printf(HELP_MSG)
	default:
		errl.Printf("wrong commanad")
		os.Exit(1)
	}
}

// init initializes few things for `ago`.
// Internally, ago uses a metadata directory for state saving. Path of the
// directory is `$HOME/.ago`. If `$HOME` is not exists, `/tmp` is used as
// default. For Android support in future, it should be `/data/local/tmp` at
// future.
//
// Hierarchy of the directory is as:
// .ago/docs/doc1
//          /info
//     /words
//
// Documents added by user resides under .ago/docs/ with its own directory. The
// document own directories are named as doc[id] which id is an integer.
// Metadata about those documents are recorded under .ago/docs/info file. The
// metadata contains original document name and current location under .ago
// directory. Because current ago support only text file, this structure is
// unnecessary. Actually, the struct is for future scaling. In future, ago will
// be an document organizer like Mendeley[1] and will support not only text
// file, but also pdf, odt, url, etc.
//
// File `words` under the .ago/ directory contains all data for words in the
// documents. It contains each word and its frequency in the docs(in total and
// per each doc), score of user, and meaning of the word. Calculated importance
// can be in there maybe but not yet decided to add it.
//
// [1] https://www.mendeley.com/
func init() {
	winfos.Wordinfos = make(map[string]wordinfo)

	if runtime.GOOS == ANDRD {
		metadat_dir = ANDRD_TMPDIR
	}

	if os.Getenv("HOME") != "" {
		metadat_dir = os.Getenv("HOME")
	}
	metadat_dir = path.Join(metadat_dir, ".ago")

	docs_dir = path.Join(metadat_dir, DOCSDIR)
	doci_path = path.Join(docs_dir, DOCINFO)
	wordi_path = path.Join(metadat_dir, WORDINFO)

	if file_exists(docs_dir) {
		read_docs_info()
		read_words_info()
		return
	}

	dbgl.Printf("docs dir is not exists. Create it.\n")
	err := os.MkdirAll(docs_dir, 0700)
	if err != nil {
		errl.Printf("docs dir %s creation failed: %s\n", docs_dir, err)
		os.Exit(1)
	}

	for _, file := range []string{doci_path, wordi_path} {
		f, err := os.Create(file)
		if err != nil {
			errl.Printf("docs info file creation failed: %s\n",
				err)
			os.Exit(1)
		}
		f.Close()

		if err = os.Chmod(file, 0600); err != nil {
			errl.Printf("chmod file %s failed: %s\n", err)
			os.Exit(1)
		}
	}

	docs_info.Next_id = NR_RSVD_DOCS
	write_docs_info()
	write_words_info()
	read_docs_info()
	read_words_info()
}

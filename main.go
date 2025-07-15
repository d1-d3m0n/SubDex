package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"os"
	"sync"
	"text/tabwriter"

	"github.com/miekg/dns"
)

type result struct {
	IPAddress string
	Hostname  string
}

func lookupA(fqdn, serveAddr string) ([]string, error) {
	var m dns.Msg
	var ips []string
	m.SetQuestion(dns.Fqdn(fqdn), dns.TypeA)
	in, err := dns.Exchange(&m, serveAddr)
	if err != nil {
		return ips, err
	}
	if len(in.Answer) < 1 {
		return ips, errors.New("no answer")
	}
	for _, answer := range in.Answer {
		if a, ok := answer.(*dns.A); ok {
			ips = append(ips, a.A.String())
		}
	}
	return ips, nil
}

func lookCNAME(fqdn, serveAddr string) ([]string, error) {
	var m dns.Msg
	var fqdns []string
	m.SetQuestion(dns.Fqdn(fqdn), dns.TypeCNAME)
	in, err := dns.Exchange(&m, serveAddr)
	if err != nil {
		return fqdns, err
	}
	if len(in.Answer) < 1 {
		return fqdns, errors.New("no answer")
	}
	for _, answer := range in.Answer {
		if c, ok := answer.(*dns.CNAME); ok {
			fqdns = append(fqdns, c.Target)
		}
	}
	return fqdns, nil
}

func lookup(fqdn, serveAddr string) []result {
	var results []result
	var cfqdn = fqdn
	for {
		cnames, err := lookCNAME(cfqdn, serveAddr)
		if err == nil && len(cnames) > 0 {
			cfqdn = cnames[0]
			continue
		}
		ips, err := lookupA(cfqdn, serveAddr)
		if err != nil {
			fmt.Println(err)
			break
		}
		for _, ip := range ips {
			results = append(results, result{IPAddress: ip, Hostname: fqdn})
		}
		break
	}
	return results
}

func worker(fqdns chan string, gather chan []result, serveAddr string, wg *sync.WaitGroup) {
	defer wg.Done()
	for fqdn := range fqdns {
		results := lookup(fqdn, serveAddr)
		if len(results) > 0 {
			gather <- results
		}
	}
}

func DisplayBanner() {
	reset := "\033[0m"
	cyan := "\033[96m"
	yellow := "\033[93m"
	green := "\033[92m"

	blue := "\033[94m"
	fmt.Printf("%s", cyan)
	fmt.Println("   _____ _     _     _____           ")
	fmt.Println("  / ____| |   | |   |   \\          ")
	fmt.Println(" | (___ | |   | | | |  | | _____  ")
	fmt.Println("  \\___ \\| |   | '_ \\| |  | |/ _ \\ \\/ /")
	fmt.Println("  ____) | |___| |_) | || |  />  < ")
	fmt.Println(" |_____/|_____|_./|_____/ \\___/_/\\_\\")
	fmt.Printf("%s", reset)
	fmt.Println()
	fmt.Printf("%s    Subdomain Enumerator Tool%s\n", yellow, reset)
	fmt.Printf("%s      made by d1_d3m0n%s\n", green, reset)
	fmt.Println()
	fmt.Printf("%s   [+] Scanning subdomains...%s\n", blue, reset)
	fmt.Printf("%s   [+] Built with Go%s\n", blue, reset)
	fmt.Printf("%s   [+] Gotta enum 'em all!%s\n", blue, reset)
	fmt.Println()
	// fmt.Println(strings().Repeat("=", 50))
	fmt.Println()
}

func main() {
	var (
		flDomain      = flag.String("domain", "", "Domain to perform enumeration.")
		flWordlist    = flag.String("wordlist", "", "Wordlist for attack.")
		flWorkerCount = flag.Int("c", 100, "The amount of workers to use.")
		flServeAddr   = flag.String("server", "8.8.8.8:53", "The DNS server to use.")
	)
	flag.Parse()
	DisplayBanner()

	if *flDomain == "" || *flWordlist == "" {
		fmt.Println("-domain and -wordlist must be provided")
		os.Exit(1)
	}

	fh, err := os.Open(*flWordlist)
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	fqdns := make(chan string, *flWorkerCount)
	gather := make(chan []result, *flWorkerCount)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < *flWorkerCount; i++ {
		wg.Add(1)
		go worker(fqdns, gather, *flServeAddr, &wg)
	}

	// Read wordlist and send FQDNs to workers
	scanner := bufio.NewScanner(fh)
	go func() {
		for scanner.Scan() {
			fqdns <- fmt.Sprintf("%s.%s", scanner.Text(), *flDomain)
		}
		close(fqdns)
	}()

	// Close gather channel when all workers are done
	go func() {
		wg.Wait()
		close(gather)
	}()

	var results []result
	for r := range gather {
		results = append(results, r...)
	}

	// Write to CSV
	fileName := *flDomain + ".csv"
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Header
	writer.Write([]string{"Hostname", "IPAddress"})

	// Output to console and file
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 4, ' ', 0)
	for _, r := range results {
		fmt.Fprintf(w, "%s\t%s\n", r.Hostname, r.IPAddress)
		if err := writer.Write([]string{r.Hostname, r.IPAddress}); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to CSV: %v\n", err)
		}
	}
	w.Flush()
	fmt.Printf("✅ Subdomains successfully saved to: %s\n", fileName)
}

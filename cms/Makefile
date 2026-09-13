run: wap lemmatize
	go test -v -run=BenchmarkCMS_1024_WAP ./... | tee run_results.txt
	rm 2600-0.txt
	rm 2600-0_lemmatized.txt

wap:
	wget https://www.gutenberg.org/files/2600/2600-0.txt

lemmatize:
	python3 lemmatize.py 2600-0.txt -o 2600-0_lemmatized.txt
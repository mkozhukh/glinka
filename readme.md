glinka
======

Yeat another Link validator

```
//collect data
glinka -t path/to/store.file http://some.com

//get stats
glinka -s path/to/store.file stats

//show errors
glinka -s path/to/store.file errors
```

You can run the command directly against the domain

```
glinka stats http://some.com
```

### Options

```
//number of concurent requests
glinka -threads 10 stats https://some.com

//log level
glinka -verbose true stats https://some.com
glinka -quiet true stats https://some.com
```

### License

MIT
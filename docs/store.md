# STORE Package 

Defines a key-value store which is used in the Redis server. 

## Architecture 
* The store is defined as a struct with a RWMutex lock and the key value store itself. 
``` 
type Store struct {
	mu   sync.RWMutex
	data map[string]string
}
``` 

* The `New()` function 


## RWMutex 
The 
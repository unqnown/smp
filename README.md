### smp

smp (**S**ay **M**y **P**assword) is a [swca1](https://github.com/unqnown/swca1) hash based CLI for strong passwords generation.

### installation

Standard `go install`:

```shell script
go install github.com/unqnown/smp
```

### configuration

To start using `smp` immediately run:

```shell script
smp init
```

Default configuration will be added to your `$HOME/.smp` directory.
You are able to override config location with `$SMPCONFIG` env variable.

```yaml
namespace: default
namespaces:
  default:
    secret: secret
    alphabet: nuls
    size: 20
    complexity: utc
```

Feel free to add more namespaces with custom options:
- secret: string of any length;
- alphabet: alphabet tokens in `wkt|alphabet` format, where `wkt` is a well-known tokens:
    * `n` or `1` - numbers;
    * `u` or `A` - uppercase letters;
    * `l` or `a` - lowercase letters;
    * `s` or `@` - symbols;
    * rest of tokens string after `|` is a custom runes which will append to result alphabet.
- size: hash size in range `[0:64]` where `0` is reserved for specifying max hash size within required complexity;
- complexity: 
    * u - ensures uniqueness of each character in hash;
    * t - guarantees any pair of adjacent characters in the hash will have a different type;
    * c - guarantees that in any pair of adjacent characters there will be no more than one letter, regardless of case.
   
### usage

To generate a new password within current namespace run:

```shell script
smp "hint"
```

To generate a new password without any configuration run:

```shell script
smp quiet -utc -s "secret" "hint"
```

For more detailed options specification run:

```shell script
smp help
````

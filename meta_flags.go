package memcache

import (
	"fmt"
	"strconv"
	"strings"
)

type metaFlager interface {
	build([]string) []string
}

func buildMetaFlags(fs []metaFlager) string {
	ss := make([]string, 0, len(fs))
	for _, f := range fs {
		ss = f.build(ss)
	}
	return strings.Join(ss, " ")
}

func obtainMetaFlagsResults(ss []string) (mr MetaResult, err error) {
	for _, f := range ss {
		k, v := f[0], f[1:]
		switch k {
		case 'W':
			mr.Won = true
		case 'Z':
			mr.Won = false
		case 'X':
			mr.Stale = true
		case 'k':
			mr.Key = v
		case 'O':
			mr.Opaque = v
		case 'c':
			mr.CasToken.value, err = strconv.ParseInt(v, 10, 64)
			mr.CasToken.setted = true
		case 'f':
			v, err := strconv.ParseUint(v, 10, 32)
			if err != nil {
				return mr, err
			}
			mr.Flags = uint32(v)
		case 'h':
			mr.Hit = (v[0] == '1')
		case 'l':
			mr.LastAccess, err = strconv.ParseUint(v, 10, 64)
		case 's':
			mr.Size, err = strconv.ParseUint(v, 10, 64)
		case 't':
			mr.TTL, err = strconv.ParseUint(v, 10, 64)
		default:
			err = fmt.Errorf("Invalid flag: %c", k)
		}
	}
	return
}

type metaFlag struct {
	k string // key
	i string // input
}

func (f metaFlag) build(of []string) []string {
	return append(of, f.k+f.i)
}

// withBinary - b: interpret key as base64 encoded binary value
func withBinary() metaFlager {
	return metaFlag{k: "b"}
}

// withCAS - c: return item cas token
func withCAS() metaFlager {
	return metaFlag{k: "c"}
}

// withFlag - f: return client flags token
func withFlag() metaFlager {
	return metaFlag{k: "f"}
}

// withHit - h: return whether item has been hit before as a 0 or 1
func withHit() metaFlager {
	return metaFlag{k: "h"}
}

// withLastAccess - l: return time since item was last accessed in seconds
func withLastAccess() metaFlager {
	return metaFlag{k: "l"}
}

// withOpaque - O(token): opaque value, consumes a token and copies back with response
func withOpaque(token string) metaFlager {
	return metaFlag{k: "O", i: token}
}

// WithQuiet - q: use noreply semantics for return codes.
func withQuiet() metaFlager {
	return metaFlag{k: "q"}
}

// withSize - s: return item size token
func withSize() metaFlager {
	return metaFlag{k: "s"}
}

// withTTL - t: return item TTL remaining in seconds (-1 for unlimited)
func withTTL() metaFlager {
	return metaFlag{k: "t"}
}

// withNoBump - u: don't bump the item in the LRU
func withNoBump() metaFlager {
	return metaFlag{k: "u"}
}

// withValue - v: return item value in <data block>
func withValue() metaFlager {
	return metaFlag{k: "v"}
}

// withVivify - N(token): vivify on miss, takes TTL as a argument
func withVivify(token uint64) metaFlager {
	return metaFlag{k: "N", i: strconv.FormatUint(token, 10)}
}

// withRecache - R(token): if token is less than remaining TTL win for recache
func withRecache(token uint64) metaFlager {
	return metaFlag{k: "R", i: strconv.FormatUint(token, 10)}
}

// withSetTTL - T(token): update remaining TTL
func withSetTTL(token uint64) metaFlager {
	return metaFlag{k: "T", i: strconv.FormatUint(token, 10)}
}

// withCompareCAS - C(token): compare CAS value when storing item
func withCompareCAS(token int64) metaFlager {
	return metaFlag{k: "C", i: strconv.FormatInt(token, 10)}
}

// withSetFlag - F(token): set client flags to token (32 bit unsigned numeric)
func withSetFlag(token uint32) metaFlager {
	return metaFlag{k: "F", i: strconv.FormatUint(uint64(token), 10)}
}

// withSetInvalid - I: invalidate. set-to-invalid if supplied CAS is older than item's CAS / - I: invalidate. mark as stale, bumps CAS.
func withSetInvalid() metaFlager {
	return metaFlag{k: "I"}
}

// withMode - M(token): mode switch to change behavior to add, replace, append, prepend
func withMode(token string) metaFlager {
	return metaFlag{k: "M", i: token}
}

// withInitialValue - J(token): initial value to use if auto created after miss (default 0)
func withInitialValue(token uint64) metaFlager {
	return metaFlag{k: "J", i: strconv.FormatUint(token, 10)}
}

// withDelta - D(token): delta to apply (decimal unsigned 64-bit number, default 1)
func withDelta(token uint64) metaFlager {
	return metaFlag{k: "D", i: strconv.FormatUint(token, 10)}
}

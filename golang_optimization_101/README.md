to avoid large value copy costs, generally, we should **NOT**:

- boxing large-size values into interfaces.
- use large-size types as function parameter types.
- use large-size types are map key and element types.

## In practice, it is encouraged to use the three-index subslice form
full slice expression: `s[low:high:max]`
it creates a similar slice to one created with `s[low:high]` except that the resulting slice’s capacity will be equal to max - low.

## Pointers in maps, arrays and channels

If the key type and element type of a map both **don’t contain pointers**, then in the scan phase of a
GC cycle, the garbage collector will not scan the entries of the map. This could save much time.

## functions inlining
After the flattening, some stack operations originally happening when calling the bar functions are saved so that code execution performance gets improved.
Inlining will make generated Go binaries larger, so compilers only inline calls to small functions

### Which functions are inline-able?
` go build -gcflags="-m -m"`
cannot inline foo: function too complex: cost 96 exceeds budget 80

**Please note that code inline costs don’t mean code execution costs.**
Since v1.18, the official standard Go compiler thinks the inline cost of for-range loop is smaller than a plain for loop.

## interfaces
From the results, we could get that boxing zero-size values, boolean values and 8-bit integer values
doesn’t make memory allocations, which is one reason why such boxing operations are much faster.
Another optimization made by the official standard Go compiler is that **no allocations are made
when boxing pointer values into interfaces**. Thus, boxing pointer values is often much faster than
boxing non-pointer values.
**The official standard Go compiler represents (the direct parts of) maps, channels and functions as
pointers internally**, so boxing such values is also as faster as boxing pointers.

A summary based on the above benchmark results (as of Go toolchain v1.18):
-  Boxing pointer values is much faster than boxing non-pointer values.
-  Boxing maps, channels and functions is as fast as boxing pointers.
-  Boxing constant values is as fast as boxing pointers.
-  Boxing zero-size values is as fast as boxing pointers.
-  Boxing boolean and 8-bit integer values is as fast as boxing pointers.
-  Boxing non-constant small values (in range [0, 255]) of any integer types (expect for 8-bit
ones) is about 3 times slower than boxing a pointer value.
-  Boxing floating-point/string/slice zero values is about 3 times slower than boxing a pointer
value.
-  Boxing a non-constant not-small integer value (out of range [0, 255]) or a non-zero
floating-point value is about (or more than) 20 times slower than boxing a pointer value.
-  Boxing non-nil slices or non-blank non-constant string values is about (or more than) 50
times slower than boxing a pointer values.
-  Boxing a struct (array) value with only one field (element) which is a small integer or a zero
bool/numeric/string/slice/point value is as faster as boxing that field (element).




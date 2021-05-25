package main

// slice operation, buf[start, stop]
// slice op gives reference to buffer memory
// reference detached if buf was reallocated (e.g. using append)
// copy(newBuf, existingBuf) checks lengths and copy only min(len1, len2) elements
// copy slice to slice, replacing existing values in newBuf

slides
======

`slides` is a simple program to create postscript slides
from simple text files. `slides` is written in [Go][go],
and requires a Go compiler (and no more) to be built using
the default build mechanism.

slides format
-------------

Each slide is described by a series of consequtive, non-empty,
lines. Lines starting with `#` are comments, and will be 
ignored, while lines with `//` are escaped lines, which will 
be printed verbatim. Each line is displayed as such, and 
will *not* be reflowed.

Comments are either just text or "command comments", designated
with an extra `+` after the `#`. These can be used to set fonts
or text sizes. All commands are listed and explained in
[this][article] introductory article to `slides`.

usage
-----

`slides` either reads it's file from standard input, or the first 
command line argument, if given. A postscript definition of all 
slides is then printed to the standard output.

Asuming `my-talk.sl` contains the contents of a presentation, 
which one would like to convert to a PostScript file, one can run

	$ slides my-talk.sl > my-talk.ps

or if one wants a PDF version of the same slides, then

	$ slides my-talk.sl | ps2pdf - my-talk.pdf

to generate these. This of course requires the `ps2pdf` programm
to be installed, from the [GhostScript][gs] package.

legal
-----

`slides` is published under CC0. See [LICENSE][license] for more
details

[go]: https://golang.org/
[gs]: https://www.ghostscript.com/
[license]: ./LICENSE
[article]: https://zge.us.to/slides.html

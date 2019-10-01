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
or text sizes. Currently, all defined comments are:

- `size`: can be set to `huge`, `large`, `normal` (default), `small` or `tiny`
- `font`: can be set to `sans` (default), `serif`, `mono`
- `style`: can be set to `bold`, `italics`, `none` (default)
- `indent`: can be set to `yes` (default), ` ` (no argument)
- `height`: can be set to an integer, an invalid number is ignored
- `width`: can be set to an integer, an invalid number is ignored
- `image`: will display an image (and only that image) on the next slide.

	Argument is a path to any image file. [ImageMagick][im] tools `convert`
	and `inspect` are required for this to work. Temporary files, of the form 
	`slides-*.eps` are created in the temporary file directory (usually 
	`/tmp`), and will *not* be automatically deleted.

Other commands are ignored, any string that is not one of the 
here listed values, is equal to the default value. 

Values are only valid for one slide. To make values valid for the 
rest of the document, or until again reset, add a `!` after the 
option name.

Values take effect at the end of a slide, i.e. an empty line. This
means that one cannot do

	#+font sans
	one line
	#+font serif
	another line

and expect the first line to be serif, and the second to be serif.

An optional colon (`:`) after the option name (including `!`) is
ignored.

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
[im]: https://imagemagick.org/index.php
[gs]: https://www.ghostscript.com/
[license]: ./LICENSE

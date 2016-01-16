# Gongram
A [nonogram](https://en.wikipedia.org/wiki/Nonogram) (aka picross, paint by numbers, picture crosswords, etc.) solver written in Go.

The solver uses the same basic technique as described on [webpbn.com](http://webpbn.com/pbnsolve.html), employing a *line solver* along with what is called *logical solving*.

It currently can solve only certain puzzles that don't require guessing. This is due to the line solver algorithm used, which is not a complete solver.
More information about the line solving algorithm can be read at webpbn.

## Usage

Clone the repository, no third party packages are required. 

Then just build it and run it as a regular go program.

    go run main.go
    
Or

    go build
    ./gongram 
    
The arguments for the program are

    Usage of ./gongram:
    -f string
        The name of the JSON file containing puzzle definitions. (default "puzzles/nonogram.json")
    -l  Displays the names in the puzzle file without solving.
    -p string
        Name of the puzzle to solve. It has to be contained in the loaded file.
        
By default, the program will display the names of the puzzles provided in the default puzzle file.
It can solve a puzzle by calling it with the `-p` argument followed by the name of the puzzle.

    ➜ ./gongram -p smiley
    Loaded puzzle: smiley
    ⎹  ×  ▉  ▉  ▉  ×  ⎸
    ⎹  ▉  ×  ▉  ×  ▉  ⎸
    ⎹  ▉  ▉  ▉  ▉  ▉  ⎸
    ⎹  ▉  ×  ×  ×  ▉  ⎸
    ⎹  ×  ▉  ▉  ▉  ×  ⎸
    
New puzzles can be loaded using the format from the included JSON file.

    {
        "name" : "smiley",
        "rows" : [[3],[1,1,1],[5],[1,1],[3]],
        "cols" : [[3],[1,1,1],[3,1],[1,1,1],[3]]
    }
    

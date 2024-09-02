# TerminalSnake
Silly snake game for terminal written in Go
It relies on [tcell](https://github.com/gdamore/tcell) as the terminal engine.

# Install
Just downlod Go and run: `go build .` inside the project directory. Then run the executable it generates.

# Obs
Some weird behavour were seen running it in different terminals. It runs better on linux terminals like Xterm.

# Controls
Use the arrow key to move the snake or the keys a-w-s-d
use p to pause and unpause. Use q to quit the game.

# Score
One point for each apple eaten. Two point for each special (green) apple eaten.

# TODO
- Create a leader board
- Create a timed special apple
- Add a menu before the game start and after the game end
- Think in more TODOS lol
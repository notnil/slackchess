# slackchess
server providing chess playing capabilities for slack

## Commands

- chess help
- chess play $user
- chess board
- chess move $notation
- chess resign
- chess draw [offer, accept, decline]

### help

The help command will show the help menu.

```
chess help
```

### play

The play command will start a game against the other user in the channel the command was sent.  There can only be one game per channel and starting a new game will end any game in progress.  

```
chess play @hoss
```

### board

The board command will show the current board.

```
chess board
```

### move

The move command takes a move in Algebraic Notation and moves the player.

```
chess move e4
```

### resign

The resign command resigns the current game.

```
chess resign
```

### draw

The draw command has three subcomands:
- offer (offers a draw to other player)
- accept (accepts a draw offer)
- decline (declines a draw offer)

The other player moving will also decline the draw offer.

```
chess draw offer
```
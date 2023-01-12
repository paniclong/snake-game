# The snake game

---

## Prerequisites
Linux only, requires ncurses library
> sudo apt-get install libncurses-dev

### Tech
**Golang 1.17**  
**C library ncurses** for read character from keyboard (without pressing enter)

### Game logic
1. `Init logger`
2. `Init game objects: snake, field`
3. `Init catch signals (SIGINT, SIGTERM, etc.) function`
4. `Init screen (ncurses)`
5. `Init pre-game screen (info)`
6. `Init loop for change direction (read arrow keys)`
7. `Init loop for stepping snake`
8. `Init main game loop (30 fps)`  
8.1 `Spawn boosters`  
8.2 `Spawn enemies`  
8.3 `Print game field` 
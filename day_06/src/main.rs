use std::{
    collections::HashSet, fs::File, io::{BufRead, BufReader}, process::exit, time
};

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
enum Direction {
    UP,
    DOWN,
    LEFT,
    RIGHT,
}

const GUARD_MARKS: [u8; 4] = [b'^', b'<', b'>', b'v'];
const OBSTACLE: u8 = b'#';

#[derive(Debug, Clone)]
struct Guard {
    position: (u16, u16),
    dir: Direction,
}

#[derive(Debug)]
struct Lab {
    dims: (u16, u16),
    obstacles: HashSet<(u16, u16)>,
    guard_start: (u16, u16),
    loops: HashSet<(u16, u16)>,
}

impl Direction {
    fn from(val: &u8) -> Self {
        match val {
            b'^' => Self::UP,
            b'<' => Self::LEFT,
            b'>' => Self::RIGHT,
            b'v' => Self::DOWN,
            _ => panic!("Invalid Direction char"),
        }
    }
}

impl Guard {
    fn turn(&mut self) {
        self.dir = match self.dir {
            Direction::UP => Direction::RIGHT,
            Direction::DOWN => Direction::LEFT,
            Direction::LEFT => Direction::UP,
            Direction::RIGHT => Direction::DOWN,
        }
    }

    fn step(&mut self) -> (u16, u16) {
        match self.dir {
            Direction::UP => self.position.1 -= 1,
            Direction::DOWN => self.position.1 += 1,
            Direction::LEFT => self.position.0 -= 1,
            Direction::RIGHT => self.position.0 += 1,
        }
        return self.position;
    }
}

impl Lab {
    fn step_guard(&mut self, guard: &mut Guard, looping: bool) -> bool {
        let next_pos: (i16, i16) = match guard.dir {
            Direction::UP => (guard.position.0 as i16, guard.position.1 as i16 - 1),
            Direction::DOWN => (guard.position.0 as i16, guard.position.1 as i16 + 1),
            Direction::LEFT => (guard.position.0 as i16 - 1, guard.position.1 as i16),
            Direction::RIGHT => (guard.position.0 as i16 + 1, guard.position.1 as i16),
        };
        if self.pos_in_dims(next_pos) {
            if self
                .obstacles
                .contains(&(next_pos.0 as u16, next_pos.1 as u16))
            {
                guard.turn();
                if looping {
                    self.search_loop(guard.clone());
                }
            }
            guard.step();
            return true;
        } else {
            return false;
        }
    }

    fn pos_in_dims(&self, pos: (i16, i16)) -> bool {
        return pos.0 < self.dims.0 as i16 && pos.1 < self.dims.1 as i16 && pos.0 >= 0 && pos.1 >= 0;
    }

    fn search_loop(&mut self, mut guard: Guard) {
        let start_pos = guard.position;
        let start_dir = guard.dir;

        while self.step_guard(&mut guard, false) {
            if guard.dir != start_dir {
                break;
            }
            if self.has_obstacle_along_slice(&guard) {      // Place ghost obstacle to attempt loop

                let mut ghost_visited = HashSet::<((u16, u16), Direction)>::new();
                let mut ghost_guard = guard.clone();
                ghost_guard.turn();
                ghost_visited.insert((ghost_guard.position, ghost_guard.dir));
                while self.step_guard(&mut ghost_guard, false) {

                    if ghost_guard.position == start_pos || ghost_visited.contains(&(ghost_guard.position, ghost_guard.dir)) {
                        let loop_obstacle = match guard.dir {
                            Direction::UP => (guard.position.0, guard.position.1 - 1),
                            Direction::DOWN => (guard.position.0, guard.position.1 + 1),
                            Direction::LEFT => (guard.position.0 - 1, guard.position.1),
                            Direction::RIGHT => (guard.position.0 + 1, guard.position.1),
                        };
                        println!("Inserting loop obstacle at {:?} with direction {:?}", loop_obstacle, guard.dir);
                        self.loops.insert(loop_obstacle);
                        break;
                    }

                    ghost_visited.insert((ghost_guard.position, ghost_guard.dir));

                }
            }
        }
    }

    fn has_obstacle_along_slice(&self, guard: &Guard) -> bool {
        return self.obstacles.iter().any(|(x, y)| match guard.dir {
            Direction::UP => *x > guard.position.0 + 1 && *y == guard.position.1,
            Direction::DOWN => *x < guard.position.0 - 1 && *y == guard.position.1,
            Direction::LEFT => *x == guard.position.0 && *y < guard.position.1 + 1,
            Direction::RIGHT => *x == guard.position.0 && *y > guard.position.1 - 1,
        });
    }
}

fn parse_input(filename: &str) -> (Lab, Guard) {
    let file = File::open(filename).unwrap();
    let reader = BufReader::new(file);

    let mut lab: Lab = Lab {
        dims: (0, 0),
        obstacles: HashSet::<(u16, u16)>::new(),
        guard_start: (0, 0),
        loops: HashSet::<(u16, u16)>::new(),
    };

    let mut guard: Guard = Guard {
        position: (0, 0),
        dir: Direction::UP,
    };

    let mut y: u16 = 0;
    for line in reader.lines().map(|l| l.unwrap()) {
        if lab.dims.0 == 0 {
            lab.dims.0 = line.len() as u16;
        }
        let line_bytes = line.as_bytes();
        for (x, tile) in line_bytes.iter().enumerate() {
            if *tile == OBSTACLE {
                lab.obstacles.insert((x as u16, y));
            } else if let Some(guard_dir) = GUARD_MARKS.iter().position(|val| val == tile) {
                guard.dir = Direction::from(&GUARD_MARKS[guard_dir]);
                guard.position = (x as u16, y);
                lab.guard_start = (x as u16, y);
            }
        }
        y += 1;
    }
    lab.dims.1 = y;
    return (lab, guard);
}

fn run_simulaton(lab: &mut Lab, guard: &mut Guard) -> HashSet<(u16, u16)> {
    let mut visited = HashSet::<(u16, u16)>::new();
    while lab.step_guard(guard, true) {
        visited.insert(guard.position);
        // println!("Visited {} nodes", visited.len());
    }
    return visited;
}

fn main() {
    let start = time::Instant::now();
    let (mut lab, mut guard) = parse_input("input");

    // println!("Parsed {:?}", lab);

    let visited_tiles = run_simulaton(&mut lab, &mut guard);

    println!("Guard visited {} tiles", visited_tiles.len());
    println!("Found {} loops", lab.loops.len());

    let (x, y) = lab.dims;
    for i in 0..y {
        for j in 0..x {
            if lab.loops.contains(&(j,i )) {
                print!("O");
            } else if visited_tiles.contains(&(j, i)) {
                print!("x");
            } else if lab.obstacles.contains(&(j, i)) {
                print!("#");
            } else  {
                print!(".")
            }
        }
        println!();
    }

    println!("Elapsed {:?}", start.elapsed())
}

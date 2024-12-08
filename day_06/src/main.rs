use std::{
    collections::HashSet,
    fs::File,
    io::{BufRead, BufReader},
    time,
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
    guard_start_dir: Direction,
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

    fn step(&mut self) {
        self.position = self.next_step();
    }

    fn next_step(&self) -> (u16, u16) {
        return match self.dir {
            Direction::UP => (self.position.0, self.position.1.wrapping_sub(1)),
            Direction::DOWN => (self.position.0, self.position.1 + 1),
            Direction::LEFT => (self.position.0.wrapping_sub(1), self.position.1),
            Direction::RIGHT => (self.position.0 + 1, self.position.1)
        }
    }
}

impl Lab {
    fn step_guard(&mut self, guard: &mut Guard, looping: bool, visited: &HashSet<(u16, u16)>) -> bool {
        let mut next_pos = guard.next_step();
        if self.pos_in_dims(next_pos) {
            while self.obstacles.contains(&next_pos) {
                guard.turn();
                next_pos = guard.next_step();
            } 
            if looping && !visited.contains(&next_pos) {
                self.search_loop(guard.clone(), next_pos, &visited);
            }
            guard.step();
            return true;
        } 
        return false;
    }

    fn pos_in_dims(&self, pos: (u16, u16)) -> bool {
        return pos.0 < self.dims.0 as u16
            && pos.1 < self.dims.1 as u16;
    }

    fn search_loop(&mut self, mut ghost_guard: Guard, obstacle: (u16, u16), visited: &HashSet<(u16, u16)>) {
        if obstacle == self.guard_start {
            return;
        }
        let orig_obstacles = self.obstacles.clone();
        self.obstacles.insert(obstacle);

        let mut ghost_visited = HashSet::<((u16, u16), Direction)>::new();

        while self.step_guard(&mut ghost_guard, false, visited) {

            if ghost_visited.contains(&(ghost_guard.position, ghost_guard.dir)) {
                if obstacle != self.guard_start {
                    self.loops.insert(obstacle);
                }
                break;
            }

            ghost_visited.insert((ghost_guard.position, ghost_guard.dir));
        }
        self.obstacles = orig_obstacles;
    }
}

fn parse_input(filename: &str) -> (Lab, Guard) {
    let file = File::open(filename).unwrap();
    let reader = BufReader::new(file);

    let mut lab: Lab = Lab {
        dims: (0, 0),
        obstacles: HashSet::<(u16, u16)>::new(),
        guard_start: (0, 0),
        guard_start_dir: Direction::UP,
        loops: HashSet::<(u16, u16)>::new(),
    };

    let mut guard: Guard = Guard {
        position: (0, 0),
        dir: Direction::UP,
    };

    let mut y: u16 = 0;
    for line in reader.lines().map(|l| l.unwrap()).filter(|l| !l.is_empty()) {
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
                lab.guard_start_dir = guard.dir;
            }
        }
        y += 1;
    }
    lab.dims.1 = y;
    return (lab, guard);
}

fn run_simulaton(lab: &mut Lab, guard: &mut Guard) -> HashSet<(u16, u16)> {
    let mut visited = HashSet::<(u16, u16)>::new();
    while lab.step_guard(guard, true, &visited) {
        visited.insert(guard.position);
    }
    return visited;
}

fn main() {
    let start = time::Instant::now();
    let (mut lab, mut guard) = parse_input("input");

    let visited_tiles = run_simulaton(&mut lab, &mut guard);

    let (x, y) = lab.dims;
    for i in 0..y {
        for j in 0..x {
            if lab.loops.contains(&(j, i)) {
                print!("O");
            } else if visited_tiles.contains(&(j, i)) {
                print!("x");
            } else if lab.obstacles.contains(&(j, i)) {
                print!("#");
            } else {
                print!(".")
            }
        }
        println!();
    }

    println!("Guard visited {} tiles", visited_tiles.len());
    println!("Found {} loops", lab.loops.len());

    println!("Elapsed {:?}", start.elapsed())
}

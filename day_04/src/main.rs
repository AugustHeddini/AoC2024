use std::{
    fs::File, 
    io::{BufRead, BufReader}
};

const SEARCH_WORD: [u8; 4] = [b'X', b'M', b'A', b'S'];
const SEARCH_CROSS_CENTER: u8 = b'A'; 
const SEARCH_CROSS: [u8; 2] = [b'M', b'S'];
const DIRECTIONS: [(i16, i16); 8] = [(1, 0), (-1, 0), (0, -1), (0, 1), (1, 1), (-1, -1), (1, -1), (-1, 1)];

fn search_for_word(word_map: &Vec<Vec<u8>>, x: usize, y: usize) -> u8 {
    let mut count: u8 = 0;

    if word_map[x][y] != SEARCH_WORD[0] {
        return count;
    }

    'dir_loop: for (dx, dy) in DIRECTIONS {
        let mut x_var = x as i16;
        let mut y_var = y as i16;

        for letter in SEARCH_WORD {
            if word_map[x_var as usize][y_var as usize] != letter {
                continue 'dir_loop;
            }
            if letter != *SEARCH_WORD.last().unwrap() {
                x_var += dx;
                y_var += dy;
                if x_var <  0 
                    || y_var < 0 
                    || x_var >= word_map.len() as i16
                    || y_var >= word_map[0].len() as i16
                {
                    continue 'dir_loop;
                }
            }
        }
        count += 1;
    }

    return count;
}

fn search_for_cross(word_map: &Vec<Vec<u8>>, x: usize, y: usize) -> u8 {
    let mut count = 0;

    if word_map[x][y] != SEARCH_CROSS_CENTER {
        return count;
    }

    if 
        (
            SEARCH_CROSS.contains(&word_map[x-1][y-1])
            && SEARCH_CROSS.contains(&word_map[x+1][y+1])
            && word_map[x-1][y-1] != word_map[x+1][y+1]
        )
        &&
        (
            SEARCH_CROSS.contains(&word_map[x-1][y+1])
            && SEARCH_CROSS.contains(&word_map[x+1][y-1])
            && word_map[x-1][y+1] != word_map[x+1][y-1]
        )
    {
        count += 1;
    }

    return count;
}

fn part1(word_map: &Vec<Vec<u8>>) {
    let mut total: u32 = 0;
    for x in 0..word_map.len() {
        for y in 0..word_map [0].len() {
            total += search_for_word(&word_map, x, y) as u32;
        }
    }

    println!("Total occurrences of {:?} found: {:?}", SEARCH_WORD, total);
}

fn part2(word_map: &Vec<Vec<u8>>) {
    let mut total: u32 = 0;
    for x in 1..word_map.len()-1 {
        for y in 1..word_map [0].len()-1 {
            total += search_for_cross(&word_map, x, y) as u32;
        }
    }

    println!("Total crosses found: {:?}", total);
}

fn parse_input(filename: &str) -> Vec<Vec<u8>> {
    let input = File::open(filename).unwrap();
    let reader = BufReader::new(input);

    let mut word_map: Vec<Vec<u8>> = vec![];

    for line in reader.lines() {
        word_map.push(line.unwrap().as_bytes().into());
    }
    
    return word_map;
}

fn main() {
    let word_map = parse_input("input");

    part1(&word_map);
    part2(&word_map);
}

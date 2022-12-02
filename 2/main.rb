$table = {
    "X" => { "value" => 1, "name" => "Rock", "points" => { "A" => 3, "B" => 0, "C" => 6 } },
    "Y" => { "value" => 2, "name" => "Paper", "points" => { "A" => 6, "B" => 3, "C" => 0 } },
    "Z" => { "value" => 3, "name" => "Scissors", "points" => { "A" => 0, "B" => 6, "C" => 3 } }
}


def help
    puts "#{$PROGRAM_NAME} input_file_path"
    exit
end

def compare(enemy, ally)
    return $table[ally]["value"] + $table[ally]["points"][enemy] 
end

def read_file(filepath)

    score = 0

    File.readlines(filepath).each do |line|
        enemy, ally = *line.split(" ")
        score += compare(enemy, ally)
    end

    puts score
end

if $PROGRAM_NAME == __FILE__

    if ARGV.length != 1 then help end

    read_file(ARGV[0])
end
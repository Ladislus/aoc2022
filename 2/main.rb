$table_1 = {
    "X" => { "value" => 1, "name" => "Rock", "points" => { "A" => 3, "B" => 0, "C" => 6 } },
    "Y" => { "value" => 2, "name" => "Paper", "points" => { "A" => 6, "B" => 3, "C" => 0 } },
    "Z" => { "value" => 3, "name" => "Scissors", "points" => { "A" => 0, "B" => 6, "C" => 3 } }
}

$table_2 = {
    "A" => { "X" => 3 + 0, "Y" => 1 + 3, "Z" => 2 + 6 },
    "B" => { "X" => 1 + 0, "Y" => 2 + 3, "Z" => 3 + 6 },
    "C" => { "X" => 2 + 0, "Y" => 3 + 3, "Z" => 1 + 6 }
}


def help
    puts "#{$PROGRAM_NAME} input_file_path"
    exit
end

def compare_1(enemy, ally)
    return $table_1[ally]["value"] + $table_1[ally]["points"][enemy] 
end

def compare_2(enemy, ally)
    return $table_2[enemy][ally]
end

def read_file(filepath)

    score_1 = 0
    score_2 = 0

    File.readlines(filepath).each do |line|
        enemy, ally = *line.split(" ")
        score_1 += compare_1(enemy, ally)
        score_2 += compare_2(enemy, ally)
    end

    puts score_1
    puts score_2
end

if $PROGRAM_NAME == __FILE__

    if ARGV.length != 1 then help end

    read_file(ARGV[0])
end
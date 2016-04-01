#!/bin/bash
( 
echo "setoption name Skill Level $1" ;
echo "position fen $2" ;
echo "go movetime 1000" ;
sleep 1
) | $3
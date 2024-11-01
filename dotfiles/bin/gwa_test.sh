#!/bin/bash

# Create a new directory
mkdir test_repo

# Change to the new directory
cd test_repo || exit

# Initialize a Git repository
git init

# Create a branch named "feature/my-module/1234/new-feature"
# git checkout -b feature/my-module/1234/new-feature

# Add a file to the branch
echo "This is a test file." >test.txt
git add test.txt
git commit -m "Add test file"

# Create a worktree
gwa "feature/my-module/YOC-1234/new-cool-feature"

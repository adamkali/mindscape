# Task Due Date Selection

## Overview

Users can optionally set a due date when creating or editing a task.

## Requirements

- User can optionally supply a due date when creating a task
- User can optionally change the due date when editing a task
- Due date is optional (nullable `due_at` column already exists in the schema)
- Due date displays in view mode (already implemented in TaskModal view section)

## Acceptance Criteria

- Creating a task without a due date works as before (null `due_at`)
- Creating a task with a due date persists the value and shows it in view mode
- Editing a task with an existing due date pre-fills the date input
- Editing a task and clearing the due date sets it to null
- The date input uses `datetime-local` format (`YYYY-MM-DDTHH:MM`)

# Structure

The repository is split into two parts: 

- The 'server' contains the backend code in go 
- The 'client' contains the frontend code with react and typescript

# Running

## Dependencies

- `go 1.23 >=`
- `node 20.14 >=`
- `npm 9.x.x >=`
- `task` (https://taskfile.dev/installation/)
- (optional) `httpie` for running scripts

## Running 

To run this whole thing locally, just run `task` from the root of the project
and it's going to setup the dependnecies and run the frontend and backend.

# Task list

- [x] I as a user can create to-do items, such as a grocery list. 
- [ ] I as another user can collaborate in real-time with user - so that we can 
(for example) edit our family shopping-list together.
- [ ] I as a user can mark to-do items as “done” - so that I can avoid clutter and focus on
things that are still pending.
- [ ] I as a user can filter the to-do list and view items that were marked as done - so that I
can retrospect on my prior progress.
- [ ] I as a user can add sub-tasks to my to-do items - so that I could make logical groups of
tasks and see their overall progress.
- [ ] I as a user can specify cost/price for a task or a subtask - so that I can track my
expenses / project cost.
- [ ] I as a user can see the sum of the subtasks aggregated in the parent task - so that in my
shopping list I can see what contributes to the overall sum. For example I can have a
task called “Salad”, where I'd add all ingredients as sub-tasks, and would see how much
a salad costs on my shopping list.
- [ ] I as a user can make infinite nested levels of subtasks.
- [ ] I as a user can add sub-descriptions of tasks in Markdown and view them as rich text
while I'm not editing the descriptions.
- [ ] I as a user can see the cursor and/or selection of another-user as he selects/types when
he is editing text - so that we can discuss focused words during our online call.
- [ ] I as a user can create multiple to-do lists where each list has its unique URL that I can
share with my friends - so that I could have separate to-do lists for my groceries and
work related tasks.
- [ ] In addition to regular to-do tasks, I as a user can add “special” typed to-do items, that will have custom style and some required fields:
    ○ ”work-task”, which has a required field “deadline” - which is a date
    ○ “food” that has fields:
    ■ required: “carbohydrate”, “fat”, “protein” (each specified in g/100g)
    ■ optional: “picture” an URL to an image used to render this item
- [ ] I as a user can keep editing the list even when I lose internet connection, and can expect it to sync up with BE as I regain connection
- [ ] I as a user can use my VR goggles to edit/browse multiple to-do lists in parallel in 3D space so that I can feel ultra-productive
- [ ] I as a user can change the order of tasks via drag & drop
- [ ] I as a user can move/convert subtasks to tasks via drag & drop
- [ ] I as a user can be sure that my todos will be persisted so that important information is not lost when server restarts
- [ ] I as an owner/creator of a certain to-do list can freeze/unfreeze a to-do list I've created to avoid other users from mutating it


# Notes

## Backend decisions

### Router choice

I am going to use the chi router, since I have used it before and it's the most
basic one that just gets out of the way and does the job.

And for defining the routes, I am going to use the OpenAPI spec and the code gen
tool. It's not really required, but I have been using it lately and I have found
that it just makes reasoning about the design a lot easier.

### Choice for ID type

I am pretty sure that the current hotness is to use UUID V7 for IDs, but I have
not really used it much yet, so I will just stick with UUID v4 that I know how
it works.

And I know that there is working google library for working with UUID V4 and
there is one for JS also. Meaning that if I want to create TODO items, they can
also be create on the frontend without having to worry about ID collisions.

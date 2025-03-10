# Preface

This todo is done as part as the homework for backend role. I didn't apply to a
fullstack role, nor the description of the role had that, BUT I know enough
react to be dangerous and made it work.

The FE took longer than I wanted to, so that is why the server part might not be
up to my standards. I also used a AI to help me with frontend side, since I have
not actually used react in quite a while, but I know enough about it that I can
make soemthing working (not all of it is good though).

The way how backend is done is not the up to my standards. We are missing greacefull shutdown,
better connection handling with channels and other things, but I sacraficed some
qulity for speed.

But I worked on this how I would have done it with a new projcet that I am
trying to prototye. 

You can mostly review the commit history, I didn't keep the greatest git commit
naming, I sacrafied it for speed.

I kept the code very simple, but extensible enough that I could keep extending
and refactoring this as the project would grow, and I wouldn't need to review
huge parts.

I included a single test just to show how I would test this. I am not a fan of
mocks, so I didn't complicate the whole thing with interfaces, just to have
mocks for testing, but for a larger project I would probably use mocks.  

Other than that the site is at:

https://ubiquitest.netlify.app

There might be some bugs that I don't know about, but for most part all the
things can be reset by logging out and refreshing the page.

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
- [x] I as another user can collaborate in real-time with user - so that we can 
(for example) edit our family shopping-list together.

- [x] I as a user can mark to-do items as “done” - so that I can avoid clutter and focus on
things that are still pending.
- [x] I as a user can filter the to-do list and view items that were marked as done - so that I
can retrospect on my prior progress.
- [x] I as a user can add sub-tasks to my to-do items - so that I could make logical groups of
tasks and see their overall progress.
- [x] I as a user can make infinite nested levels of subtasks.


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

### Skip auth for most routes

I am just trying to move fast here and I am not going to worry about security
here. 

In real scenarios, I would add oauth and authentication middlewars to the routes (or even
just have a gateway infront of the server that does that, for example KrakenD
for multiservice deployments).

Then there would also be logic of autherization, if the user even can view the
other user's todos and info. Probabaly based on organizations where we would
have N-N relationship for users and organizations. This would allow user with
one login to be part of multiple organizations.

This would heavily also impact the DB design, which is the reason why I am not
adding that right now. We would have to add org_id to each task, and it just
complicates things more than I have time for it.

### Choice for DB

It's just a sqlite db that lives on the fly.io container that is running. I
didn't just have the time to setup turso for sqlite or some other managed db
service right now.

### Choice for WebSockets

I am using websockets for realtimes updates. They allow for bidirectional data
transfer and it's great for this use case.


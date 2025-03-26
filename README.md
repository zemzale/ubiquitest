# Ubiquiti Test

Deployed frontend at: https://ubiquitest.netlify.app

# Preface

This todo app is done as part of the homework for backend role. I didn't apply to a
fullstack role, nor the description of the role had any metions of React/TS. But I know enough
React to be dangerous so I made it work.

The FE took longer than I wanted to, so that is why the server part might not be
up to my standards. I used quite a lot AI to help me with frontend side, since I have
not actually used react in quite a while, but I was able to get it working.

There are some thigns missing on the backend side that I would usually include
like:
- greacefull shutdown,
- proper authentication
- better logging
- tracing
- and other things

But I had very little time to do this, so I sacraficed them for speed, not
because I don't know how to do them. This is a good example of a prototye I
would create, not a production ready app.

You can mostly review the commit history, I didn't keep the greatest git commit
naming, also sacrafied it for speed.

I kept the code very simple, but extensible enough that I could keep extending
and refactoring this as the project would grow, and I wouldn't need to review
huge parts.

I used some DDD by trying to model the domain entities and their actions, but
there are some parts of users and maybe tasks that are not fully switched over
to the DDD approach yet. I also inlined all the database logic since splitting
it out for project of this size would have been a waste. Usually that would be
part of another module.

I included a single test just to show how I would test this. I am not a fan of
mocks, so I didn't complicate the whole thing with interfaces, just to have
mocks for testing, but for a larger project I would probably use mocks.  

Use a unique name for the login, and you can create tasks and see others create
them in real time in one board.

There might be some bugs that I don't know about, but for most part all the
things can be reset by logging out and refreshing the page.

--- 

### EDIT

I have now managed to seperate the storage layer from the business logic. They
now exist in domain usecases that are seperated from the storage layer.
This is the pattern I have started to use
lately since it allows for easy composition and isolated testing. The oteher way
is by creating this functionality in "services" that group similar functionality,
but that makes it much more coupled.

Before you ask about why I don't use interfaces. I do know what they are, and
how to use them, BUT:
- I am not writing tests that would require mocking
- Mocking is an anti-pattern IMO
- There are no multiple implementations of the same thing, so abstracting now
makes nose sense. In a long term project, I would create the wrong abstraction
here and that would be more painful than it is now. 

I follow the principle of "don't make the decisions you don't need to make". An
interface is a decision that you make when you need to abstract implementation
details, and try to stick to it since you will have to change multiple things,
if you break the abstraction.

If I am not creating multiple implementations of the same thing, I don't need
to make that choice, so I just skip it.

I managed to slip in a the feature of creating tasks with costs and calculating
data for that.

As for monitoring:

- And added basic prometheus metrics. 
- I though of setting up Sentry for tracing and error reporting, but I don't see a
real point in showing that I know how to read the sentry documentation.
- And I won't setup open telemetry for tracing and metrics, unless someone is
paying me. Those SDKs are garbage.

Other than that, I could just implement more features and add more tests, but I
with the time I have, and seeing that they won't really change much of showing
off my skills, I am not going to do that.

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
- [x] I as a user can specify cost/price for a task or a subtask - so that I can track my
expenses / project cost.


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

Then there would also be logic of authorization, if the user even can view the
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


### Choice for DI

I used the `samber/lo` library for dependency injection. It's a really easy and
nice library to use and it just let's me avoid writing some boilerplate and
passing things around that much. 

Most of the projcets I create end up with something that is very similar, except
I write it by hand. Have been using it for a while now, and it's quite nice.

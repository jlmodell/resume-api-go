# resume-api-go
simple api built using golang and the framework go-fiber (v2) to be consumed by a simple front end

# routes
[x] GET / = basic api information  
[x] GET /api/resume = resume  

[x] PUT /api/resume/skills = add a skill  
[x] PUT /api/resume/projects = add a project  
[x] PUT /api/resume/certifications = add a certificate  
[x] PUT /api/resume/links = add a link  

[x] DELETE /api/resume/:skills|projects|certifications|links> = delete a specific item from one of the arrays in the database  
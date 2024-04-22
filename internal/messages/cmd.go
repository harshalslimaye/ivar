package messages

const InitCmdWelcome = `
%s init walks you through creating a %s file with the essentials, suggesting handy defaults. 
For more info, type 'ivar help init'. Remember, you can press ^C anytime to quit. 
Afterward, use 'ivar install <pkg>' to add dependencies. Easy peasy!

`
const InitCmdAlreadyExists = `
Oops! Found a package.json file here already. 
If you want to start fresh, just delete it and 
run 'ivar init' again.
`

While developing the connector, please fill out this form. This information is needed to write docs and to help other users set up the connector.

## Connector capabilities

1. What resources does the connector sync?
    This connector syncs:
        — Users
        — Roles

2. Can the connector provision any resources? If so, which ones?
   The connector can provision:
        — Roles for Users

## Connector credentials

1. What credentials or information are needed to set up the connector? (For example, API key, client ID and secret, domain, etc.)
   This connector requires an API Key and a Domain. Args: --bearer-token and --domain

2. For each item in the list above:

    * How does a user create or look up that credential or info? Please include links to (non-gated) documentation, screenshots (of the UI or of gated docs), or a video of the process.
        1- Log in, then in the top right corner of the main page of your fluid topics page, click on administration.
        2- In the menu that opens, click on integrations.
        3- On the integrations page, below the list of api keys, a section will appear to create an api-key by adding a name and clicking create&add.
        4- When you click on create&add it will open a menu of options to customize the apikey. 
        5- After configuring your api key click on ok, and when the menu closes, in the integrations page with the apikey list click on save in the lower right corner.
   
        Note: documentation of api keys: [Fluid-topics-APIKEY](https://doc.fluidtopics.com/r/Fluid-Topics-Configuration-and-Administration-Guide/Configure-a-Fluid-Topics-tenant/Integrations/API-keys)
        
    * Does the credential need any specific scopes or permissions? If so, list them here.
      To validate the use of the connector, the apikey must be configured to have the ADMIN role, since it allows synchronization.
      You can configure it by creating the new apikey (step four of the instructions above) or by editing it by clicking on the edit button which is like a pencil.
        
    * If applicable: Is the list of scopes or permissions different to sync (read) versus provision (read-write)? If so, list the difference here.
      To be able to do both, the ADMIN role is required.
   
    * What level of access or permissions does the user need in order to create the credentials? (For example, must be a super administrator, must have access to the admin console, etc.)  
      The user must have the role of ADMIN to be able to access the page and create the credential 
{ pkgs, ... }: {
  channel = "stable-23.11"; # or "unstable"

  # Use https://search.nixos.org/packages to find packages
  packages = [
    pkgs.go
  ];

  # Sets environment variables in the workspace
  env = {
    API_KEY = "AIzaSyAw-UQBZuKm2e-r3PjgvBudbVTOZLMd3nU";
  };
  idx = {
    # Search for the extensions you want on https://open-vsx.org/ and use "publisher.id"
    extensions = [
      "golang.go"
    ];

    # Enable previews
    previews = {
      enable = true;
      previews = [
        # {
        #   # Example: run "npm run dev" with PORT set to IDX's defined port for previews,
        #   # and show it in IDX's web preview panel
        #   command = ["npm" "run" "dev"];
        #   manager = "web";
        #   id = "web";
        #   env = {
        #     # Environment variables to set for your server
        #     PORT = "$PORT";
        #   };
        # }
      ];
    };

    # Workspace lifecycle hooks
    workspace = {
      # Runs when a workspace is first created
      onCreate = {
         default.openFiles = [
          "adventure.go" "README.md"
         ];
        # Example: install JS dependencies from NPM
        # npm-install = 'npm install';
      };
      onStart = {
        # Example: start a background task to watch and re-build backend code
        # watch-backend = "npm run watch-backend";
      };
    };
  };
}

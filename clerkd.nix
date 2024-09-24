{ clerk }:
clerk.overrideAttrs {
  pname = "clerkd";
  subPackages = [ "cmd/clerkd" ];

  postInstall = "";
  meta.mainProgram = "clerkd";
}

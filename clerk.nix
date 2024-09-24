{ buildGoModule, installShellFiles }:
buildGoModule {
  pname = "clerk";
  version = "0.1.0";
  src = ./.;
  vendorHash = "sha256-qK9xeNIdfs0IIL7U9f7qDgIv8s4mCuOKyPDqX+H55sU=";
  doCheck = false;

  subPackages = [ "cmd/clerk" ];

  nativeBuildInputs = [ installShellFiles ];

  postInstall = ''
    installShellCompletion --cmd clerk \
      --bash <($out/bin/clerk completion bash) \
      --zsh <($out/bin/clerk completion zsh) \
      --fish <($out/bin/clerk completion fish)
  '';

  meta.mainProgram = "clerk";
}

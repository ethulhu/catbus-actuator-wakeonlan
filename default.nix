# SPDX-FileCopyrightText: 2020 Ethel Morgan
#
# SPDX-License-Identifier: MIT

{ pkgs ? import <nixpkgs> {} }:
with pkgs;

buildGoModule rec {
  name = "catbus-wakeonlan-${version}";
  version = "latest";
  goPackagePath = "go.eth.moe/catbus-wakeonlan";

  modSha256 = "19nr5pw4c4jap051djg4m1j83p3vkfirhvjvw2ld6mk3zjblbh2f";

  src = ./.;

  meta = {
    homepage = "https://ethulhu.co.uk/catbus";
    licence = stdenv.lib.licenses.mit;
  };
}

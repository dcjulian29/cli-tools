trap [System.Exception] {
  "Exception: {0}" -f $_.Exception.Message
  return
}

$baseDir = (Resolve-Path $("$PSScriptRoot")).Path
$artifactDir = Join-Path -Path $baseDir -ChildPath "dist"

if (Test-Path $artifactDir) {
  Remove-Item -Path $artifactDir -Recurse -Force
}

New-Item -Path $artifactDir -ItemType Directory | Out-Null

Get-ChildItem -Path $baseDir -Directory | ForEach-Object {
  Push-Location $(Join-Path -Path $baseDir -ChildPath $_.Name)

  if ((Test-Path ./.goreleaser.yml) -or (Test-Path ./.goreleaser.yaml)) {
    if (Get-Command -Name goreleaser -ErrorAction SilentlyContinue) {
      Write-Output "`n---------------------->  $pwd"
      goreleaser build --single-target --clean --snapshot
      Write-Output "`n`n"
    }
  }

  if (Test-Path dist/artifacts.json) {
    $dist = Get-Content -Path dist/artifacts.json -Raw | ConvertFrom-Json

    foreach ($item in $dist) {
      if ($item.type -eq "Binary") {
        $src = Join-Path -Path $PWD -ChildPath $item.path
        $dst = Join-Path -Path $artifactDir -ChildPath $item.name

        Write-Output "$src`n  |`n  V`n$dst"

        Move-Item -Path $src -Destination $dst
      }
    }

    Remove-Item -Path dist/ -Recurse -Force
  }

  Pop-Location
}

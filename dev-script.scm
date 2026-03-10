(use-modules (guix profiles)
             (guix packages)
             (gnu packages docker)
             (guix build utils))

(system* "docker" "run"
         "-d"
         "--name" "postgres-dev"
         "-e" "POSTGRES_USER=postgres"
         "-e" "POSTGRES_PASSWORD=postgres"
         "-p" "5432:5432"
         "-v" "postgres-data:/var/lib/postgresql/data"
         "postgres:16")

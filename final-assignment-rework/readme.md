## Final Assignment - Voting Application

This folder contains my final assignment for CST680 - Cloud Native Software Engineering

To initialize my project, first run ./build-all-docker.sh then you can start and stop using docker compose up/down

Once the containers are running, you can run ./loadcache.sh to show all posts work (you may want to clear cache-data that is a part of my repository as it has loadcache already stored - you can do this by running ./clearcache.sh). Then you can run the same script a second time to see that all posts should fail since they were already added. All data persists through stopping and starting of the container, so once loadcache.sh is run once, the cache will never be empty again unless the cache-data is erased manually.

After, you can run ./getcachegood.sh to show GET successes for existing data stored. Then run ./getcachebad.sh to show GET failures for ID's that do not exist in the cache or empty locations (such as a voter who does not have any votes in their vote history).
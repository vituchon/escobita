// Routines for synchronizing running instance data with local data
package services

// Dev notes: THe sync methods should be independent of filepaths (local) or urls (remotes), maybe an struct can be build at
// app boot time using proper config file.
// Example:
// dbSyncStateStore = services.NewDbSyncStateStore(config), config refers to a config intrisic to services (not to model!)
// dbSyncStateStore.Load() error
// dbSyncStateStore.Save() error

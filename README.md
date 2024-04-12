# Filestores

![filestores](gopher_small.png "Filestores")

I was getting tired creating implementations to save files in local or cloud. By now local and aws is supported.

## How to use

Create a struct which implements `filestores.Storable`:

```go
type Storable interface {
    Filename() string
    GetStoreInfo() ObjectInfo
}
```

**IMPORTANT:** When you implement `Filename() string` method remember not to put trailing slash `/` at the beginning of the name and without file extension. It should be a return like `my/path/to/file`. This may cause unintended results

And select one of the available two options:

### Local

Uses filesystem to save files. To use it you can:

```golang
func main() {
    filestore := filestores.NewFileSystem("./uploads")
    //...
}
```

### AWS

Uses AWS S3 API to save files. It is important to know this package uses AWS Go SDK V2 and your `filestores.Storable.Filename` method return a valid S3 URL.

```golang
func main() {
    filestore := filestores.NewAwsBucket("aBucketNameYouWant", aws.Config{
        //...
    })
    //...
}
```

## Helpers

There are two helpers to create a `filestores.ObjectInfo` instance need to implement `filestores.Storable.GetStoreInfo`. These are:

- `filestores.FileToStoreInfo`
- `filestores.FileHeaderToStoreInfo`

## Reminders

This just cover the basic need to save files...by now.

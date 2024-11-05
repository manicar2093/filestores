# Filestores

**🚧 This library as under development until release of version `1.x.x` 🚧**

![filestores](gopher_small.png "Filestores")

I was getting tired creating implementations to save files in local or cloud. By now local and aws is supported.

## Motivation

I'm not a big fan of pay when I'm developing something new. This is where comes `filestores` idea. 

This I mainly use it when in develop to use localhost to serve files and, when I deploy my project, I can configure an AWS Bucket (only option available by now), without changing many code.

## How to use

Create a struct which implements `filestores.Storable`:

```go
type Storable interface {
    Filename() string
    GetStoreInfo() ObjectInfo
}
```

Like this:

```go
type StorablePicture struct {
	Id uuid.UUID
	FileHeader           *multipart.FileHeader
}

func (c StorablePicture) Filename() string {
	return fmt.Sprintf("%s", c.Id)
}

func (c StorablePicture) GetStoreInfo() filestores.ObjectInfo {
	info, _ := filestores.FileHeaderToStoreInfo(c.FileHeader)
	return info
}
```

**IMPORTANT:** When you implement `Filename() string` method remember not to put trailing slash `/` at the beginning of the name and without file extension. It should be a return like `my/path/to/file`. This may cause unintended results

And select one of the available two options:

### Local

Uses filesystem to save files. To use it you can:

```golang
func main() {
    filestore := filestores.NewFileSystem("./uploads", "http://localhost:5000")
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

## What comes next?

I found [go-cloud](https://github.com/google/go-cloud) project which aims to do possible to use any cloud implementation using just one library. I'm thinking to implement it with `filestores` to reach all go-cloud supported storage services.

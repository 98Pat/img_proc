
# Image Processing Collection by Patrick Protte

Welcome to the **Image Processing Collection**, a versatile and efficient tool for applying various filters to images using Go. This project is my first step into the world of Go, and serves as a learning project to explore the capabilities of the language while providing a useful tool for image processing.

## Features

- **Multi-processor support**: Utilize the power of multiple logical processors for faster image processing.
- **Multiple filters**: Apply various filters such as blur, invert, comic, spot, edge, heat, and gaussian blur.
- **Customizable options**: Each filter comes with its own set of configurable parameters to fine-tune the output.

## Installation

To install the Image Processing Collection, ensure you have Go installed on your system. Then, clone the repository and build the project:

```bash
git clone https://github.com/98Pat/img_proc.git
cd img_proc
make
```

## Usage

The program can be executed from the command line with various flags to control its behavior. Below is the general syntax and a description of the available options:

```bash
./img_proc-[linux | windows] [options]
```

### Options

Non-flag arguments cannot be followed by flag-arguments.

- `-I int`  
  **Description**: Number of iterations the filter should be applied.  
  **Default**: 1  

- `-c int`  
  **Description**: Number of logical processors to use.  
  **Default**: Maximum available  

- `-f string`  
  **Description**: Type of filter to apply.  
  **Required Arguments**: Depends on the filter type.  
  **Filter Types**:
  - `blur`
  - `invert`
  - `comic`        (optional: color step count (int), default 3)
  - `spot`         (required: posX, posY, radius (int, int, float))
  - `edge`         (optional: amplification (int), default 1)
  - `heat`
  - `gaussianblur` (optional: kernel size/radius, sigma (int, float), default 5, 2.0)

- `-h`  
  **Description**: Display help information.

- `-i string`  
  **Description**: Path to the input image file.  
  **Required**

- `-o string`  
  **Description**: Path to the output image file.  
  **Required**

### Example Usage

#### Blur Filter

Apply a blur filter to an image:

```bash
./img_proc-linux -i input.jpg -o output.jpg -f blur
```

#### Comic Filter

Apply a comic filter with 5 color steps:

```bash
./img_proc-linux -i input.jpg -o output.jpg -f comic 5
```

#### Spot Filter

Apply a spot filter with specific position and radius:

```bash
./img_proc-linux -i input.jpg -o output.jpg -f spot 100 150 50.0
```

#### Edge Filter

Apply an edge filter with amplification of 2:

```bash
./img_proc-linux -i input.jpg -o output.jpg -f edge 2
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

**Patrick Protte**

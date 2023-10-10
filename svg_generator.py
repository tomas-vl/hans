#!/usr/bin/env python3
import defopt
import drawsvg as svg
import imagesize
from dataclasses import dataclass
from copy import deepcopy


@dataclass
class Picture:
    """Ahoj."""
    picture_filepath: str
    picture: svg.Drawing
    width: int
    height: int
    border_size: int

    def true_width(self) -> int:
        return self.width + 2 * self.border_size

    def true_height(self) -> int:
        return self.height + 2 * self.border_size

def InitializePicture(width: int, height: int, border_size: int, picture_filepath: str) -> Picture:
    svg_picture: svg.Drawing = svg.Drawing(
        width + 2 * border_size, 
        height + 2 * border_size, 
        origin=(0, 0), 
        font_family='Libertinus Serif'
    )
    wrapped_picture: Picture = Picture(picture_filepath, svg_picture, width, height, border_size)
    return wrapped_picture

def createBaseSVG(owp: Picture, frequency: str, duration: str) -> Picture:
    # owp: Picture = deepcopy(iwp)
    owp.picture.append(
        svg.Image(
            owp.border_size, 
            owp.border_size, 
            owp.width, 
            owp.height, 
            owp.picture_filepath, 
            embed=True
        )
    )
    owp.picture.append(
        svg.Rectangle(
            owp.border_size,
            owp.border_size, 
            owp.width, owp.height, 
            fill='none', stroke='black', stroke_width=3
        )
    )
    owp.picture.append(
        svg.Text(
            f'Frekvence 0—{frequency} Hz', 
            font_size=owp.border_size - 15, 
            x=owp.border_size - 10, 
            y=owp.true_height() - owp.border_size, 
            transform=f'rotate(-90,{owp.border_size - 10},{owp.true_height() - owp.border_size})'
        )
    )
    owp.picture.append(
        svg.Text(f'Trvání {duration} s', 
            font_size=owp.border_size - 15,
            x=owp.width + owp.border_size, 
            y=owp.border_size - 10, 
            text_anchor='end'
        )
    )
    return owp

def addLine(owp: Picture, position: int) -> Picture:
    # owp: Picture = deepcopy(iwp)
    
    owp.picture.append(
        svg.Line(
            position, 
            owp.border_size, 
            position, 
            owp.height + owp.border_size, 
            stroke='red', 
            stroke_width=3, 
            stroke_dasharray='5 3'
        )
    )
    return owp;

def addLettersAroundPosition(owp: Picture, position: int, letter_1: str, letter_2: str) -> Picture:
    position_shift: int = 15
    owp.picture.append(
        svg.Text(
            f'{letter_1}', 
            font_size = owp.border_size - 15, 
            x = position - position_shift, 
            y = owp.height + owp.border_size + 25, 
            text_anchor = 'end'
        )
    )
    owp.picture.append(
        svg.Text(
            f' {letter_2}', 
            font_size = owp.border_size - 15, 
            x = position + position_shift, 
            y = owp.height + owp.border_size + 25, 
            text_anchor = 'start'
        )
    )
    return owp

def addLetterBetweenPositions(owp: Picture, left_pos: int, right_pos: int, letter: str) -> Picture:
    letter_pos: int = (left_pos + right_pos) // 2
    owp.picture.append(
        svg.Text(
            f'{letter}', 
            font_size = owp.border_size - 15, 
            x = letter_pos, 
            y = owp.height + owp.border_size + 25, 
            text_anchor = 'end'
        )
    )
    return owp

def main():
    width: int = 0
    height: int = 0
    width, height = [int(x) for x in imagesize.get('test_input.png')]
    print(width, height)

    output_picture: Picture = InitializePicture(width, height, 40, 'test_input.png')
    output_picture = createBaseSVG(output_picture, '1848', '18,48')

    output_picture = addLine(output_picture, 350)
    output_picture = addLettersAroundPosition(output_picture, 350, 'a', 'é')

    output_picture = addLine(output_picture, 740)
    output_picture = addLine(output_picture, 900)
    output_picture = addLetterBetweenPositions(output_picture, 740, 900, 'ř')
    output_picture.picture.save_svg('example.svg')

if __name__ == '__main__':
	defopt.run(main)

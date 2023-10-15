#!/usr/bin/env python3
'''
    File name: svg_generator.py
    Author: Tomáš Vlček <tvlcek@mail.muni.cz>
    Date created: 2023-10-10
    License: GNU General Public License v3.0
    Python Version: ≥3.11.5
'''

import defopt
import drawsvg as svg
import imagesize
from dataclasses import dataclass
from copy import deepcopy


@dataclass
class Picture:
    '''Wrapper around the SVG canvas.'''
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
        font_family='Charis SIL'
    )
    wrapped_picture: Picture = Picture(picture_filepath, svg_picture, width, height, border_size)
    return wrapped_picture

def createBaseSVG(owp: Picture, frequency: str, duration: str) -> Picture:
    '''Fills the SVG canvas with the bitmap, border and margin, and labels with frequency and duration.'''
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
            f'Frekvence 0–{frequency} Hz', 
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

def createBareBaseSVG(owp: Picture) -> Picture:
    '''Fills the SVG canvas with the bitmap, border and margin. No labels are added.'''
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
    return owp

def createFinalSVG(iwp: Picture, lines: list[int], letters: list[tuple[int, str]], frequency: str, duration: str):
    '''Fills the picture with everything, lines and all.'''
    owp: Picture = createBaseSVG(iwp, frequency, duration)
    for line_pos in lines:
        owp = addLine(owp, line_pos)

    for pos_and_letter in letters:
        owp = addLetterAtPosition(owp, pos_and_letter[0], pos_and_letter[1])

    return owp

def addLine(owp: Picture, position: int) -> Picture:
    '''Adds a vertical line to the SVG canvas at the specified `position`.'''
    # owp: Picture = deepcopy(iwp)
    
    owp.picture.append(
        svg.Line(
            position, 
            owp.border_size, 
            position, 
            owp.height + owp.border_size, 
            stroke='red', 
            stroke_width=3, 
            stroke_dasharray='6 12'
        )
    )
    return owp;

def addLetterAtPosition(owp: Picture, position: int, letter: str) -> Picture:
    owp.picture.append(
        svg.Text(
            f'{letter}', 
            font_size = owp.border_size - 15, 
            x = position, 
            y = owp.height + owp.border_size + 25, 
            text_anchor = 'middle',
        )
    )
    return owp

def addLettersAroundPosition(owp: Picture, position: int, letter_1: str, letter_2: str) -> Picture:
    '''Adds `letter_1` and `letter_2` to the bottom margin of the SVG canvas around the specified `position`.'''
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
    '''Adds `letter` to the bottom margin of the SVG canvas between `left_pos` and `right_pos`.'''
    letter_pos: int = (left_pos + right_pos) // 2
    owp.picture.append(
        svg.Text(
            f'{letter}', 
            font_size = owp.border_size - 15, 
            x = letter_pos, 
            y = owp.height + owp.border_size + 25, 
            text_anchor = 'middle',
        )
    )
    return owp


def main():
    width: int = 0
    height: int = 0
    width, height = [int(x) for x in imagesize.get('test_input/01.png')]
    print(width, height)

    output_picture: Picture = InitializePicture(width, height, 40, 'test_input/01.png')
    lines = [30, 150, 350, 600]
    letters = [
        (90, 'á'),
        (250, 'óó'),
        (475, 'další písmeno')
    ]
    output_picture = createFinalSVG(output_picture, lines, letters, '50', '500')

    output_picture.picture.save_svg('output.svg')

if __name__ == '__main__':
	defopt.run(main)

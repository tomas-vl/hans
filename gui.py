#!/usr/bin/env python3
'''
    File name: gui.py
    Author: Tomáš Vlček <tvlcek@mail.muni.cz>
    Date created: 2023-10-10
    License: GNU General Public License v3.0
    Python Version: ≥3.11.5
'''

from cairosvg import svg
import defopt
import PySimpleGUI as sg
import imagesize
from dataclasses import dataclass
import svg_generator

@dataclass
class UserInput:
    filepath: str = ''
    frequency: str = ''
    duration: str = ''

def readjustLetterPositions(line_positions: list[int], letter_positions: list[tuple[int, str]], start: int, end: int) -> list[tuple[int, str]]:
    temp_line_positions: list[int] = line_positions.copy()
    temp_line_positions.append(start)
    temp_line_positions.append(end)
    temp_line_positions = sorted(temp_line_positions)

    updated_letter_positions: list[tuple[int, str]] = letter_positions.copy()

    for i in range(len(letter_positions)):
        num: int = letter_positions[i][0]
        
        for j in range(len(temp_line_positions) - 1):
            if temp_line_positions[j] <= num <= temp_line_positions[j + 1]:
                average = (temp_line_positions[j] + temp_line_positions[j + 1]) // 2
                updated_letter_positions[i] = (average, letter_positions[i][1])
                break

    return updated_letter_positions

def main():
    user_input: UserInput = UserInput()
    line_positions: list[int] = list()
    letter_positions: list[tuple[int, str]] = list()
    temp_spectrogram: svg_generator.Picture = None

    col_labels = [
        [sg.Text('Frekvence:')],
        [sg.Text('Trvání:')],
    ]
    col_inputs = [
        [sg.Input(enable_events=True, key='-FREQ-')],
        [sg.Input(enable_events=True, key='-DURATION-')],
    ]
    col_graph_controls = [
        [sg.R('Čára', 1, key='-LINE-', enable_events=True)],
        [sg.R('Písmeno', 1, key='-LETTER-', enable_events=True)],
       # [sg.R('Smazat prvek', 1, key='-ERASE-', enable_events=True)],
    ]

    layout = [
        [
            sg.Input(key='-INPUT-'),
            sg.FileBrowse(file_types=(('PNG Images', '*.png'), ('ALL Files', '*.*'))),
            sg.Button('Load file', key='-OPEN_FILE-'),
            sg.Push()
        ],
        [
            sg.Col(col_labels), 
            sg.Col(col_inputs),
        ],
        [
            sg.Graph((0, 0), (0, 0), (0, 0), key='-GRAPH-', enable_events=True, change_submits=True, drag_submits=False),
            sg.Col(col_graph_controls),
        ],
        [
            sg.Push(),
            sg.Button('Save file', key='-SAVE_FILE-')
        ],
    ]

    window = sg.Window('hans', layout, resizable=True)

    image = window['-GRAPH-']
    if isinstance(image, sg.ErrorElement):
        return 1

    while True:
        event, values = window.read()
        print(event, values)
        print((
            f'fr:\t{user_input.frequency}\n'
            f'tr:\t{user_input.duration}\n'
            f'input:\t{user_input.filepath}\n'
            f'lines:\t{line_positions}\n'
            f'letters:\t{letter_positions}\n'
        ))

        if event == sg.WINDOW_CLOSED or event == 'Quit':
            break

        elif event == '-OPEN_FILE-':
            # načtení cesty k neupravené bitmapě spektrogramu
            user_input.filepath = values['-INPUT-']
            print(f'Otvírám soubor {user_input.filepath}')


            new_width, new_height = [int(x) for x in imagesize.get(user_input.filepath)]   # načtení rozlišení původní bitmapy

            temp_spectrogram = svg_generator.InitializePicture(new_width, new_height, 40, user_input.filepath)
            temp_spectrogram = svg_generator.createBareBaseSVG(temp_spectrogram)
            temp_spectrogram.picture.save_png('temp.png')
            
            image.set_size((temp_spectrogram.true_width(), temp_spectrogram.true_height()))
            image.change_coordinates((0,0), (temp_spectrogram.true_width(), temp_spectrogram.true_height()))

            window['-GRAPH-'].erase()
            image.draw_image('temp.png', location=(0,image.CanvasSize[1]))

        elif event == '-SAVE_FILE-':
            if user_input.frequency == '' or user_input.duration == '' or user_input.filepath == '':
                sg.popup('Chyba: Nebyla zadána jedna z požadovaných hodnot.', title='Chyba!', grab_anywhere=True, keep_on_top=True, any_key_closes=True)
            else:
                output_filepath: str = sg.popup_get_file("Save As", save_as=True)
                print(f'Ukládám soubor {user_input.filepath} do {output_filepath}')

                # vygenerování finálního obrázku
                output_spectrogram = svg_generator.InitializePicture(temp_spectrogram.width, temp_spectrogram.height, 40, user_input.filepath)
                output_letter_positions: list[tuple[int, str]] = readjustLetterPositions(line_positions, letter_positions, output_spectrogram.border_size, output_spectrogram.border_size + output_spectrogram.width)
                output_spectrogram = svg_generator.createFinalSVG(output_spectrogram, line_positions, output_letter_positions, user_input.frequency, user_input.duration)
                output_spectrogram.picture.save_svg(output_filepath)

        elif event == '-FREQ-':
            user_input.frequency = values['-FREQ-']
            print(f'načtena frekvence {user_input.frequency}')

        elif event == '-DURATION-':
            user_input.duration = values['-DURATION-']
            print(f'načteno trvání {user_input.duration}')
        elif event == '-GRAPH-':    # kliknutí do náhledu
            if values['-LINE-']:    # v menu je zvolena 'čára', pročež se bude kreslit čára
                x: int = 0
                y: int = 0
                x, y = values['-GRAPH-']    # získání polohy kurzoru
                image.draw_line(            # vykreslení čáry do náhledu
                    (x,temp_spectrogram.border_size),   # souřadnice náhledu začínají v levém dolním rohu
                    (x,temp_spectrogram.border_size + temp_spectrogram.height), 
                    width=4, 
                    color='red'
                )
                line_positions.append(x)    # uložení x-ové souřadnice čáry pro finální vykreslení
            elif values['-LETTER-']:
                x: int = 0
                y: int = 0
                x, y = values['-GRAPH-']
                input_letter = sg.popup_get_text("Vepiš písmeno:") # získání písmena od uživatele 
                if input_letter != None:
                    image.draw_text(input_letter, location=(x, 10))
                    letter_positions.append((x, input_letter))  # uložení x-ové souřadnice písmena pro finální vykreslení
            # elif values['-ERASE-']:
            #     for figure in drag_figures:
            #         image.delete_figure(figure)
            # elif values['-CLEAR-']:
            #     image.erase()

    window.close()

    return

if __name__ == '__main__':
	defopt.run(main)
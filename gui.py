#!/usr/bin/env python3
'''
    File name: gui.py
    Author: Tomáš Vlček
    Date created: 2023-10-10
    License: GNU General Public License v3.0
    Python Version: ≥3.11.5
'''

import defopt
import PySimpleGUI as sg

def main():
    layout = [
        [
            sg.Input(key='-INPUT-'),
            sg.FileBrowse(file_types=(("PNG Images", "*.png"), ("ALL Files", "*.*"))),
            sg.Button('Load file', key='-OPEN_FILE-'),
            sg.Push()
        ],
        [sg.Text('Frekvence:'), sg.Input(enable_events=True, key='-FREQ-')],
        [sg.Text('Trvání:'), sg.Input(enable_events=True, key='-DURATION-')],
        [sg.Graph((800, 800), (0, 450), (450, 0), key='-GRAPH-', enable_events=True, change_submits=True, drag_submits=False)],
        [
            sg.Push(),
            sg.Button('Save file', key='-SAVE_FILE-')
        ],
    ]

    window = sg.Window('hans', layout, resizable=True)

    input_filename: str = str()

    frequency: str = str()
    duration: str = str()

    while True:
        event, values = window.read()
        print(event, values)
        print((
            f'fr:\t{frequency}\n'
            f'tr:\t{duration}\n'
            f'input:\t{input_filename}\n'
        ))

        if event == sg.WINDOW_CLOSED or event == 'Quit':
            break
        elif event == '-OPEN_FILE-':
            input_filename = values['-INPUT-']
            print(f'Otvírám soubor {input_filename}')
        elif event == '-SAVE_FILE-':
            print(f'Ukládám soubor {input_filename}')
            # Získání složky kam pomocí sg.popup_get_folder
            # uložení pomocí output_picture.picture.save_svg('output.svg')
        elif event == '-FREQ-':
            frequency = values['-FREQ-']
            print(f'načtena frekvence {frequency}')
        elif event == '-DURATION-':
            duration = values['-DURATION-']
            print(f'načteno trvání {duration}')            

     # Finish up by removing from the screen
    window.close()

    return

if __name__ == '__main__':
	defopt.run(main)
from pathlib import Path

import pandas as pd
import yaml


def convert(input_file: Path, output_dir: Path):
    data = read_xlsx(input_file)
    write_yaml(data, output_dir / (input_file.stem + ".yaml"))


def read_xlsx(file: Path) -> dict:
    result = {}
    xl = pd.ExcelFile(file)
    for sheet_name in xl.sheet_names[0:2]:
        result[sheet_name] = {}
        df = xl.parse(sheet_name, index_col=0)
        df.dropna(inplace=True)
        value_column = df.columns[0]
        result[sheet_name]["unit"] = value_column
        result[sheet_name]["data"] = {
            date.date().isoformat(): values[value_column]
            for date, values in df.to_dict(orient="index").items()
        }
    return result


def write_yaml(data: dict, file: Path) -> None:
    if not file.parent.exists():
        file.parent.mkdir(parents=True)
    file.write_text(yaml.dump(data))

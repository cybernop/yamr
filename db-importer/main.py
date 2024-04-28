import argparse
from dataclasses import dataclass
import logging
from pathlib import Path
from typing import Dict, List, Tuple

import pandas as pd
import requests


logger = logging.getLogger(__name__)


@dataclass
class DataConfig:
    name: str
    unit: str
    data: List[Tuple]


def read_excel(file: str, sheet_name: str, column_name: str):
    content = pd.read_excel(Path(file), sheet_name=sheet_name, index_col=0)
    content.dropna(inplace=True)
    data = content[column_name].to_dict()
    data = [(date.isoformat(), entry) for date, entry in data.items()]
    return data


def read_data(file: str, data_configs: List[Dict]) -> List[DataConfig]:
    result = []
    for config in data_configs:
        sheet_data = read_excel(file, config["sheet"], config["column"])
        result.append(DataConfig(config["name"], config["unit"], sheet_data))
    return result


def post_data(data: List[DataConfig], host: str):
    kinds = get_kinds(host)
    for dataset in data:
        kind_id = get_kind_id(dataset.name, kinds)
        if kind_id is None:
            logger.info("creating new kind in DB for %s", dataset.name)
            try:
                create_kind(host, dataset.name, dataset.unit)
            except Exception as e:
                logger.error("failed to create new kind: %s", e)
                exit(1)

            kinds = get_kinds(host)
            kind_id = get_kind_id(dataset.name, kinds)

        for recorded_on, reading in dataset.data:
            try:
                create_reading(host, kind_id, recorded_on, reading)
            except Exception as e:
                logger.error("failed to create reading: %s", e)
                continue


def get_kinds(host: str):
    result = requests.get(f"http://{host}/kind")
    json_content = result.json()
    return json_content["kinds"]


def create_kind(host: str, name: str, unit: str):
    data = {"name": name, "unit": unit}
    result = requests.post(f"http://{host}/kind", json=data)
    if result.status_code != 201:
        raise Exception(f"could not create kind '{name}': {result.text}")


def get_kind_id(kind: str, kinds: List[str]) -> int | None:
    for k in kinds:
        if k["name"] == kind:
            return k["id"]
    return None


def create_reading(host: str, kind_id: int, recorded_on: str, reading):
    data = {"kind_id": kind_id, "recordedOn": f"{recorded_on}Z", "reading": reading}
    result = requests.post(f"http://{host}/reading", json=data)
    if result.status_code != 201:
        raise Exception(
            f"could not create reading '{kind_id, recorded_on,reading}': {result.text}"
        )


def get_args():
    parser = argparse.ArgumentParser()
    parser.add_argument("--config", type=Path, required=True)
    return parser.parse_args()


if __name__ == "__main__":
    import yaml

    args = get_args()
    config = yaml.safe_load(args.config.read_text())
    data = read_data(config["file"], config["data"])
    post_data(data, config["host"])

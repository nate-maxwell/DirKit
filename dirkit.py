"""
# Dir Kit

* A simple toolkit for folder and file handling that eliminates
boilerplate or wraps commonly used functions in a consistent
namespace for easy rememberance/importing.
"""


import os
import re
import datetime
import shutil
import json
import platform
from typing import Optional
from typing import Union
from pathlib import Path


SAFETY_PATH = Path('D:/safety/')  # Change on per-project needs.


def get_dir_contents(path: Path, full_path: bool = False) -> Union[list[Path], list[str], None]:
    """
    Gets the content names, or full path for contents, of a directory.

    Args:
        path (pathlib.Path): Directory path to list contents of.

        full_path (bool): To return string names or paths of directory contents. Defaults to False.

    Returns:
        String names of directory contents if full_path = False,
        Paths for directory contents if full_path = True.
    """
    if path.exists():
        if full_path:
            contents = list(path.glob('*'))
        else:
            contents = os.listdir(path)

        return contents

    return None


def create_directory(path: Path) -> Path:
    """
    Creates a directory from the given path.

    Args:
        path (path): The directory path for the directory to create.

    Returns:
        str: The created, or pre-existing, directory path.
    """
    if not os.path.exists(path):
        os.makedirs(path)

    return path


def create_dated_directory(path: Path) -> Path:
    """
    Creates a directory with today's date as the name.

    Args:
        path (str): The path, with base directory name, to place the directory.

    Returns:
        str: The full path of the created directory with date.
    """
    date_path = Path(path, get_date())
    create_directory(date_path)

    return date_path


def delete_safe_directory(path: Path) -> None:
    """
    Deletes a directory and its contents as long as they are within the safety path.

    Args:
        path (Path): the path to the directory to delete. Will throw exception if
        path is not within the safety path.
    """
    if SAFETY_PATH in list(path.parents):
        for root, dirs, files in os.walk(path, topdown=False):
            for name in files:
                filename = os.path.join(root, name)
                os.remove(filename)
            for name in dirs:
                os.rmdir(os.path.join(root, name))

        os.rmdir(path)
    else:
        raise ValueError(f'Path must be within {SAFETY_PATH.as_posix()}!')


def delete_safe_file(filepath: Path) -> None:
    """
    Removes specified file as long as it is within the safety path.

    Args:
        filepath (Path): the path to the file to delete. Will throw exception if
        path does not start with the safety path.
    """
    if SAFETY_PATH in list(filepath.parents):
        os.remove(filepath)
    else:
        raise ValueError(f'Path must be within {SAFETY_PATH.as_posix()}!')


def delete_safe_files_in_directory(directory_path: Path) -> None:
    """
    Delete all files in a directory as long as they are within the safety path.

    Args:
        directory_path (Path): Path to the directory.
    """
    try:
        files = get_dir_contents(directory_path, True)
        for file in files:
            delete_safe_file(filepath=file)
        print("All files deleted successfully.")
    except OSError:
        print("Error occurred while deleting files.")


def copy_file(source: Path, destination: Path, new_name: Optional[str] = '') -> None:
    """
    Copy file into a separate destination folder.

    Args:
        source (Path): file path of the file to copy.

        destination (Path): folder path of where to copy the file to.

        new_name (Optional[str]): an optional argument to rename the file.
    """
    if source.parent == destination:
        if source.is_dir():
            return
        else:
            if '.' in new_name and not new_name.split('.')[-1].isnumeric():
                new_base_name = os.path.splitext(new_name)[0]
            else:
                new_base_name = new_name

            ext = source.name.split('.')[-1]
            replace_name = f'{new_base_name}.{ext}'
            shutil.copy(source, Path(source.parent, replace_name))
    else:
        create_directory(destination)

        if new_name:
            target = Path(destination, new_name)
        else:
            target = destination

        shutil.copy(source, target)


def copy_folder_contents(source: Path, destination: Path) -> None:
    """
    Copy contents of a folder to the given destination.

    Args:
        source (Path): folder path to the folder that is to be copied.

        destination (Path): folder path to copy the folder + contents to.
    """
    shutil.copytree(source, destination, dirs_exist_ok=True)


def get_date() -> str:
    """Returns str: 'YYYYMMDD'"""
    today = datetime.date.today()
    return today.strftime("%Y%m%d")


def get_time() -> str:
    """Returns str: 'HH:MM:SS:XX', X is microsecond."""
    now = datetime.datetime.now().time().isoformat()[:-4]
    return now


def get_os_info() -> tuple[str, str, str]:
    """Returns tuple[str, str, str]: OS name, release number, and version number."""
    system  = platform.system()
    release = platform.release()
    version = platform.version()
    return system, release, version


def export_data_to_json(path: Path, data: dict, overwrite: bool = False) -> None:
    """
    Export dict to json file path.

    Args:
        path (Path): the file path to place the .json file.

        data (dict|list): the data to export into the .json file.

        overwrite(bool): to overwrite json file if it already exists in path.
            Defaults to False.
    """
    if not path.exists() or overwrite:
        with open(path, 'w') as outfile:
            json.dump(data, outfile, indent=4)
    else:
        return


def import_data_from_json(filepath: Path) -> Optional[dict]:
    """
    Import data from a .json file.

    Args:
        filepath (Path): the filepath to the json file to extract data from.

    Returns:
        any: will return data if json file exists, None if it doesn't.
    """
    if os.path.exists(filepath):
        with open(filepath) as file:
            data = json.load(file)
            return data

    return None


def sort_path_list(path_objs: list[Path] = None) -> Optional[list[Path]]:
    """
    Alpha-numerically sorts a list of pathlib.Paths.

    Args:
        path_objs(list[Path]): The list of paths to sort.

    Returns:
        Optional[list[Path]]: The sorted list of paths or None if no
        list was provided.
    """
    if path_objs is None:
        return None

    if len(path_objs) == 1:
        return path_objs

    sort_strings = []
    for p in path_objs:
        sort_strings.append(p.as_posix())

    convert = lambda text: int(text) if text.isdigit() else text
    alphanum_key = lambda key: [convert(c) for c in re.split('([0-9]+)', key)]
    sort_strings.sort(key=alphanum_key)

    sorted_paths = []
    for s in sort_strings:
        sorted_paths.append(Path(s))

    return sorted_paths

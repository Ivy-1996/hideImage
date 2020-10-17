import cv2
import numpy as np


def encode_image(show_image_path, hide_image_path, out_image_name):
    '''自动计算代码'''
    show_image = cv2.imread(show_image_path)

    hide_image = cv2.imread(hide_image_path)

    width, height = show_image.shape[:2]

    hide_image = cv2.resize(hide_image, (height, width), interpolation=cv2.INTER_CUBIC)

    show_image = show_image // 10 * 10

    hide_image = hide_image * (9 / 255)

    out_image = show_image + hide_image

    cv2.imwrite(out_image_name, out_image)

    '''
    # 手动计算代码
    
    show_image = cv2.imread(show_image_path)

    hide_image = cv2.imread(hide_image_path)

    shape = show_image.shape

    show_image_data = show_image.data.tolist()

    hide_image_data = hide_image.tolist()

    out_image = np.zeros((shape[0], shape[1], 3), dtype=np.uint8)

    for x in range(shape[0]):
        for y in range(shape[1]):
            show_point = show_image_data[x][y]
            hide_point = hide_image_data[x][y]
            new_show_r, new_show_g, new_show_b = [value // 10 * 10 for value in show_point]
            new_hide_r, new_hide_g, new_hide_b = [value * 9 / 255 for value in hide_point]

            out_point = new_hide_r + new_show_r, new_hide_g + new_show_g, new_show_b + new_hide_b
            out_image[x][y] = out_point

    cv2.imwrite(out_image_name, out_image)
    '''


def decode_image(target_image_path, show_image_name, hide_image_name):
    '''自动计算代码'''
    target_image = cv2.imread(target_image_path)

    show_image = target_image // 10 * 10

    hide_image = (target_image - show_image) * (255 / 9)

    cv2.imwrite(show_image_name, show_image)

    cv2.imwrite(hide_image_name, hide_image)

    '''
    # 手动计算代码
    target_image = cv2.imread(target_image_path)

    target_image_data = target_image.data.tolist()

    shape = target_image.shape

    show_image = np.zeros((shape[0], shape[1], 3), dtype=np.uint8)

    hide_image = np.zeros((shape[0], shape[1], 3), dtype=np.uint8)

    for x in range(shape[0]):
        for y in range(shape[1]):
            target_point = target_image_data[x][y]

            show_point = [value // 10 * 10 for value in target_point]
            hide_point = [(target - show_point[target_point.index(target)]) * 255 / 9 for target in
                          target_point]
            show_image[x][y] = show_point
            hide_image[x][y] = hide_point

    cv2.imwrite(show_image_name, show_image)
    cv2.imwrite(hide_image_name, hide_image)
    '''


if __name__ == '__main__':
    encode_image("a.png", "b.png", "c.png")
    decode_image("c.png", "d.png", "e.png")

#!/usr/bin/env python3

import numpy as np

mode = 'default'

if mode == 'favicon':
    size = 64
    levels = [5, 12]
    stroke_width = 2
elif mode == 'small':
    size = 128
    levels = [5, 7, 12]
    stroke_width = 1
elif mode == 'default' or True:
    size = 512
    levels = [5, 7, 10, 15]
    stroke_width = 2

output = 'logo.svg'
growth_factor = 0.16
base = 5
steps = 50
# for golden spiral: base = ln(phi)/(pi/2)
colors = ['#4ab5c4', '#4ebcb3', '#4dbfa2', '#3797db', '#47add8', '#59c4cd', '#62d9a7', '#46b6a9']
bg = '#fff'
stroke = '#000'

border = size * 0.05

pi = np.pi

def spiral_r(theta):
    return np.exp(growth_factor * theta)

def main():
    rng = np.random.RandomState(1)
    previous_color = None
    def random_color():
        nonlocal previous_color
        color = previous_color
        while color == previous_color:
            color = rng.choice(colors)
        previous_color = color
        return color

    half = size/2

    lines = []
    strokes = []
    lines.append('<?xml version="1.0" encoding="UTF-8"?>')
    lines.append('<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="%d" height="%d">' % (size, size))
    # lines.append('<rect style="fill: %s; stroke: none;" width="%d" height="%d" />' % (bg, size, size))


    pts = []
    spaces = []
    for i, v in enumerate(levels):
        spaces.append(np.linspace(2*pi*i, 2*pi*(i+1), v))
    rep = len(spaces)
    width = spiral_r(2*pi*(rep-1)+3*pi/2) - spiral_r(2*pi*(rep-1)+pi/2)
    height = spiral_r(2*pi*rep) - spiral_r(2*pi*(rep-1)+pi)
    # XXX the math here is a bit sketchy
    scale = (half - border) / (2 * max(width, height))
    x0 = (border + scale*spiral_r(2*pi*(rep-1)+3*pi/2))
    y0 = (border + scale*spiral_r(2*pi*rep))
    def point(theta):
        r = scale * spiral_r(theta)
        x = x0 - r * np.cos(theta+pi/2)
        y = y0 - r * np.sin(theta+pi/2)
        return x, y

    # quads or triangles in inner part
    inside = spaces[0]
    curr = 1
    while True:
        if curr + 2 < len(inside):
            # quad
            index = [0, curr, curr+1, curr+2]
        elif curr + 1 < len(inside):
            index = [0, curr, curr+1]
        else:
            break
        thetas = [inside[i] for i in index]
        points = [
            '%f,%f' % point(t)
            for t in thetas
        ]
        fill = random_color()
        lines.append('<polygon points="%s" style="fill: %s; stroke: none;" />' % (' '.join(points), fill))
        strokes.append('<polyline points="%s" style="fill: none; stroke: %s; stroke-width: %d" />' % (' '.join(points[-1:] + points[:1]), stroke, stroke_width))
        curr = curr+2

    # polys in outer parts
    for level_in_i in range(len(spaces)-1):
        level_in = spaces[level_in_i]
        level_out = spaces[level_in_i+1]
        prev = 0
        for pt_i in range(len(level_in)-1):
            pt = level_in[pt_i]
            pt_n = level_in[pt_i+1]
            # XXX logic here could use some work probably
            closest = prev
            n_closest = np.argmin(np.abs((level_out - 2*pi) - pt_n))
            if n_closest == closest:
                n_closest += 1
            pt_closest = level_out[closest]
            pt_n_closest = level_out[n_closest]
            prev = n_closest
            # draw quad
            outer = range(n_closest, closest-1, -1)
            outer_t = [level_out[i] for i in outer]
            points = [
                '%f,%f' % point(t)
                for t in [pt, pt_n] + outer_t
            ]
            fill = random_color()
            lines.append('<polygon points="%s" style="fill: %s; stroke: none;" />' % (' '.join(points), fill))
            strokes.append('<polyline points="%s" style="fill: none; stroke: %s; stroke-width: %d" />' % (' '.join(points[-1:] + points[:1]), stroke, stroke_width))

    # outline
    space = np.concatenate(spaces)
    for theta in space:
        x, y = point(theta)
        pts.append('%f,%f' % (x, y))
    pts.append('%f,%f' % point(spaces[-1][0]))
    strokes.append('<polyline points="%s" style="fill: none; stroke: %s; stroke-width: %d" />' %
            (' '.join(pts), stroke, stroke_width))


    lines.extend(strokes)
    lines.append('</svg>')
    with open(output, 'w') as f:
        f.write('\n'.join(lines))
    print('ok')

if __name__ == '__main__':
    main()

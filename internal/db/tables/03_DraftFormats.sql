CREATE TABLE IF NOT EXISTS draft_formats (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    cube_url VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Insert default format for Vintage Cube
INSERT INTO draft_formats (name, description, cube_url)
VALUES ('PowerMack', 'Pmax without bad cards', 'https://cubecobra.com/cube/list/PowerMack')
ON DUPLICATE KEY UPDATE id=id;

INSERT INTO draft_formats (name, description, cube_url)
VALUES ('LSVCube', "It's a classic Vintage Cube, with an eye towards powerful gameplay, and underperforming cards get rotated out frequently. Enjoy!", 'https://cubecobra.com/cube/list/LSVCube')
ON DUPLICATE KEY UPDATE id=id; 
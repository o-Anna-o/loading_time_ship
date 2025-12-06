// Prevents additional console window on Windows in release, DO NOT REMOVE!!
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]
use tauri_plugin_http;
fn main() {
    tauri::Builder::default()
        .plugin(tauri_plugin_http::init())
        .plugin(tauri_plugin_cors_fetch::init())
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}

// #![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

// fn main() {
//    src_tauri_lib::run()
//}
